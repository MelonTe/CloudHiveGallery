package user

//用户更新请求
type UserEditRequest struct {
	ID          uint64 `json:"id,string" swaggertype:"string"` //用户ID
	UserName    string `json:"userName"`                       //用户昵称
	UserProfile string `json:"userProfile"`                    //用户简介
}
