package picture

//返回图片的基本标签和分类
type PictureTagCategory struct {
	TagList      []string `json:"tagList"`
	CategoryList []string `json:"categoryList"`
}
