package websocket

import (
	"chg/internal/common"
	"chg/internal/consts"
	"chg/internal/ecode"
	"chg/internal/manager/websocket/model/request"
	"chg/internal/manager/websocket/model/response"
	"chg/internal/model/entity"
	resUser "chg/internal/model/response/user"
	"chg/internal/service"
	"encoding/json"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/gorilla/websocket"
	"log"
	"net/http"
	"strconv"
	"sync"
)

var upgrader = websocket.Upgrader{
	CheckOrigin: func(r *http.Request) bool {
		// 可以加入更严格的域名验证
		return true
	},
}

//处理图片编辑请求消息

// 处理升级请求，需要先校验权限
var userService = service.NewUserService()
var pictureService = service.NewPictureService()
var spaceService = service.NewSpaceService()

// 定义会话session
type PictureEditSessions struct {
	Sessions    sync.Map // pictureId -> *SessionBucket 存储处在当前图片编辑会话里的角色session
	EditingUser sync.Map //pictureId -> userID 存储当前正在编辑图片的角色
}
type SessionBucket struct {
	sync.Mutex
	Clients []*PictureEditClient
}
type PictureEditClient struct {
	user *entity.User
	conn *websocket.Conn
}

var sessionManager = &PictureEditSessions{}

// 获取当前Picture的会话桶，不存在则返回一个新的
func (p *PictureEditSessions) GetOrCreateBucket(pictureId uint64) *SessionBucket {
	val, _ := p.Sessions.LoadOrStore(pictureId, &SessionBucket{})
	return val.(*SessionBucket)
}

// 添加一个会话
func (p *PictureEditSessions) AddClient(pictureId uint64, client *PictureEditClient) {
	bucket := p.GetOrCreateBucket(pictureId)
	bucket.Lock()
	defer bucket.Unlock()
	bucket.Clients = append(bucket.Clients, client)
}

// 退出会话
func (p *PictureEditSessions) RemoveClient(pictureId uint64, client *PictureEditClient) {
	val, ok := p.Sessions.Load(pictureId)
	if !ok {
		return
	}
	bucket := val.(*SessionBucket)
	//注意，当前对象可能正在持有EditingUser的锁，所以需要在这里解锁
	//使用该方法会获取锁，所以要在bucket获取锁之前调用
	HandleExitAction(nil, client.user, pictureId, client)
	bucket.Lock()
	//在服务器中移除掉该会话的session
	for i, c := range bucket.Clients {
		if c == client {
			// 移除
			bucket.Clients = append(bucket.Clients[:i], bucket.Clients[i+1:]...)
			break
		}
	}
	// 如果没剩下人了，可以删掉这个 bucket
	if len(bucket.Clients) == 0 {
		p.Sessions.Delete(pictureId)
	}
	// 发送广播，用户退出编辑状态
	editResponse := &response.PictureEditResponseMessage{
		Type:    consts.WS_PICTURE_EDIT_MESSAGE_EXIT_EDIT,
		Message: "用户 " + client.user.UserName + " 退出编辑",
		User:    resUser.GetUserVO(*client.user),
	}
	//广播之前释放锁
	bucket.Unlock()
	BoardCastToPicture(pictureId, editResponse, nil)
}

// 协议升级入口
func PictureEditHandShake(c *gin.Context) {
	type WSPictureEditRequest struct {
		PictureID uint64 `form:"pictureId"` // 图片ID
	}
	wsPictureEditRequest := &WSPictureEditRequest{}
	var originErr error
	wsPictureEditRequest.PictureID, originErr = strconv.ParseUint(c.Query("pictureId"), 10, 64)
	if originErr != nil {
		common.BaseResponse(c, nil, "请求体解析失败", ecode.PARAMS_ERROR)
		return
	}
	if wsPictureEditRequest.PictureID == 0 {
		common.BaseResponse(c, nil, "图片ID不能为空", ecode.PARAMS_ERROR)
		return
	}
	loginUser, err := userService.GetLoginUser(c)
	if err != nil {
		common.BaseResponse(c, nil, err.Msg, err.Code)
		return
	}
	//校验用户是否有编辑图片的权限
	picture, err := pictureService.GetPictureById(wsPictureEditRequest.PictureID)
	if err != nil {
		common.BaseResponse(c, nil, err.Msg, err.Code)
		return
	}
	var space *entity.Space
	//检查若图片是属于空间的，则需要检查用户是否有权限
	if picture.SpaceID != 0 {
		space, err = spaceService.GetSpaceById(picture.SpaceID)
		if err != nil {
			common.BaseResponse(c, nil, err.Msg, err.Code)
			return
		}
		//检查用户是否有权限
		if space.SpaceType != consts.SPACE_TEAM {
			//不是团队空间，拒绝握手
			common.BaseResponse(c, nil, "非团队空间，拒绝握手", ecode.PARAMS_ERROR)
			return
		}
	}
	//若图片是属于公共空间的，为了兼容未来扩展，允许管理员协同编辑，现在需要校验是否用户拥有编辑权限即可
	permissionList := service.GetPermissionList(space, loginUser)
	haveEditPermission := false
	for _, v := range permissionList {
		if v == "picture:edit" {
			haveEditPermission = true
		}
	}
	if !haveEditPermission {
		common.BaseResponse(c, nil, "无权限进行协同操作", ecode.NO_AUTH_ERROR)
		return
	}
	//权限校验完毕，升级协议
	conn, originErr := upgrader.Upgrade(c.Writer, c.Request, nil)
	if originErr != nil {
		log.Println("WebSocket 升级失败:", err)
		common.BaseResponse(c, nil, "WebSocket 升级失败", ecode.SYSTEM_ERROR)
		return
	}
	//协议升级成功，广播连接成功的消息
	editResponse := &response.PictureEditResponseMessage{
		Type:    consts.WS_PICTURE_EDIT_MESSAGE_INFO,
		Message: "用户 " + loginUser.UserName + " 进入编辑状态",
		User:    resUser.GetUserVO(*loginUser),
	}
	BoardCastToPicture(picture.ID, editResponse, nil)
	//记录当前用户的会话
	client := &PictureEditClient{
		user: loginUser,
		conn: conn,
	}
	sessionManager.AddClient(picture.ID, client)
	//进入后续websocket处理消息逻辑
	go WSPictureEditHandler(client, loginUser, picture.ID)
}

// 广播方法，通知处在当前图片编辑会话的所有角色，发生了操作改变
// 广播方法会获取锁
// excludeClient: 被排除广播外的客户端
func BoardCastToPicture(pictureId uint64, editResponse *response.PictureEditResponseMessage, excludeClient *PictureEditClient) {
	bucket := sessionManager.GetOrCreateBucket(pictureId)
	bucket.Lock()
	defer bucket.Unlock()
	data, _ := json.Marshal(editResponse)
	//便利所有会话
	for _, client := range bucket.Clients {
		if client == excludeClient {
			continue // 排除当前客户端
		}
		err := client.conn.WriteMessage(websocket.TextMessage, data)
		if err != nil {
			log.Println("WebSocket 广播失败:", err)
			client.conn.Close()
			continue
		}
	}
}

// 收到前端发送的消息，进行消息处理
func TextMessageHandler(curClient *PictureEditClient, loginUser *entity.User, pictureId uint64, text []byte) {
	//解析消息体
	reqText := &request.PictureEditRequestMessage{}
	if err := json.Unmarshal(text, reqText); err != nil {
		//解析异常
		log.Println("解析消息体失败:", err)
		return
	}
	//对消息类型进行处理
	switch reqText.Type {
	case consts.WS_PICTURE_EDIT_MESSAGE_EDIT_ACTION:
		//处理编辑动作
		HandleEditAction(reqText, loginUser, pictureId, curClient)
	case consts.WS_PICTURE_EDIT_MESSAGE_ENTER_EDIT:
		//进入编辑状态
		HandleEnterAction(reqText, loginUser, pictureId, curClient)
	case consts.WS_PICTURE_EDIT_MESSAGE_EXIT_EDIT:
		//退出编辑状态
		HandleExitAction(reqText, loginUser, pictureId, curClient)
	default:
		//编辑消息错误相应
		editResponse := &response.PictureEditResponseMessage{
			Type:    consts.WS_PICTURE_EDIT_MESSAGE_ERROR,
			Message: "未知消息类型",
			User:    resUser.GetUserVO(*loginUser),
		}
		//只广播当前的前端
		data, _ := json.Marshal(editResponse)
		curClient.conn.WriteMessage(websocket.TextMessage, data)
	}
}

// 进入编辑操作
func HandleEnterAction(reqText *request.PictureEditRequestMessage, loginUser *entity.User, pictureId uint64, curClient *PictureEditClient) {
	//若当前用户无人编辑，才进行处理
	if _, ok := sessionManager.EditingUser.Load(pictureId); !ok {
		//设置当前用户为编辑者
		sessionManager.EditingUser.Store(pictureId, loginUser.ID)
		//广播所有人，有人加入了编辑
		resMsg := &response.PictureEditResponseMessage{
			Type:    consts.WS_PICTURE_EDIT_MESSAGE_ENTER_EDIT,
			Message: "用户 " + loginUser.UserName + " 开始编辑图片",
			User:    resUser.GetUserVO(*loginUser),
		}
		BoardCastToPicture(pictureId, resMsg, nil)
	}
}

// 处理编辑动作
func HandleEditAction(reqText *request.PictureEditRequestMessage, loginUser *entity.User, pictureId uint64, curClient *PictureEditClient) {
	//执行编辑操作时，首先需要判断当前用户是否是编辑者
	if editer, ok := sessionManager.EditingUser.Load(pictureId); ok {
		editerId := editer.(uint64)
		if editerId == loginUser.ID {
			//是编辑者，执行操作

			//检查操作是否存在
			if !consts.IsEditAction(reqText.EditAction) {
				//操作不存在，打印日志并返回即可
				log.Println("编辑操作不存在:", reqText.EditAction)
				return
			}
			//构造响应，发送具体的操作通知
			resMsg := &response.PictureEditResponseMessage{
				Type:       consts.WS_PICTURE_EDIT_MESSAGE_EDIT_ACTION,
				Message:    fmt.Sprintf("用户 %s 执行了编辑操作: %s", loginUser.UserName, consts.GetActionName(reqText.EditAction)),
				User:       resUser.GetUserVO(*loginUser),
				EditAction: reqText.EditAction,
			}
			//广播给出了自己以外的所有人，否则造成重复编辑
			BoardCastToPicture(pictureId, resMsg, curClient)
		}
	}
}

// 退出编辑操作
func HandleExitAction(reqText *request.PictureEditRequestMessage, loginUser *entity.User, pictureId uint64, curClient *PictureEditClient) {
	//校验当前用户是否是编辑者
	if editer, ok := sessionManager.EditingUser.Load(pictureId); ok {
		editerId := editer.(uint64)
		if editerId == loginUser.ID {
			//是编辑者，执行退出操作
			sessionManager.EditingUser.Delete(pictureId)
			//广播所有人，有人退出了编辑
			resMsg := &response.PictureEditResponseMessage{
				Type:    consts.WS_PICTURE_EDIT_MESSAGE_EXIT_EDIT,
				Message: "用户 " + loginUser.UserName + " 退出编辑图片",
				User:    resUser.GetUserVO(*loginUser),
			}
			BoardCastToPicture(pictureId, resMsg, nil)
		}
	}
}

// 定义任务
type MessageTask struct {
	Client    *PictureEditClient
	User      *entity.User
	PictureId uint64
	message   []byte
}

// 创建协程+channel队列处理任务
func WSPictureEditHandler(curClient *PictureEditClient, loginUser *entity.User, pictureId uint64) {
	//保持连接中
	defer sessionManager.RemoveClient(pictureId, curClient)
	taskChan := make(chan MessageTask, 10)
	//为当前conn创建一个守护协程，按顺序处理消息而不堵塞
	go func() {
		for task := range taskChan {
			//处理消息
			TextMessageHandler(task.Client, task.User, task.PictureId, task.message)
		}
	}()
	for {
		//尝试获取数据
		_, msg, err := curClient.conn.ReadMessage()
		if err != nil {
			//断开连接
			log.Printf("用户 %s 断开连接", loginUser.UserName)
			break
		}
		//发送消息给守护协程即可返回
		taskChan <- MessageTask{
			Client:    curClient,
			User:      loginUser,
			PictureId: pictureId,
			message:   msg,
		}
	}
}
