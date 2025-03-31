<template>
  <div class="picture-upload">
    <a-upload
      list-type="picture-card"
      :show-upload-list="false"
      :custom-request="handleUpload"
      :before-upload="beforeUpload"
    >
      <img v-if="picture?.url" :src="picture?.url" alt="avatar" />
      <div v-else>
        <loading-outlined v-if="loading"></loading-outlined>
        <plus-outlined v-else></plus-outlined>
        <div class="ant-upload-text">点击或拖拽上传图片</div>
      </div>
    </a-upload>
  </div>

</template>
<script lang="ts" setup>
import { ref } from 'vue';
import { PlusOutlined, LoadingOutlined } from '@ant-design/icons-vue';
import { message } from 'ant-design-vue';
import type { UploadChangeParam, UploadProps } from 'ant-design-vue';
import { postPictureUpload } from '@/api/picture.ts'

interface Props {
  picture?: API.PictureVO
  spaceId?: string
  onSuccess?: (newPicture: API.PictureVO) => void
}

const props = defineProps<Props>();
// 上传图片
// file：上传的文件
const handleUpload = async ({file} : any) => {
  loading.value = true;
  try {
    const params: API.PictureUploadRequest = props.picture ? {id : props.picture.id} : {}
    params.spaceId = props.spaceId;
    const res = await postPictureUpload(params,file,{} )
    if (res.data.code === 0 && res.data.data){
      message.success("图片上传成功")
      //上传成功的信息传递给副组件
      props.onSuccess?.(res.data.data);
    } else{
      message.error("图片上传失败，", + res.data.msg);
    }
  }catch (error) {
    console.log("图片上传失败",error)
    message.error("图片上传失败，"+ error.message);
  }
  loading.value = false;
}
const loading = ref<boolean>(false);

// 前端校验
const beforeUpload = (file: UploadProps['fileList'][number]) => {
  //校验图片格式
  const isJpgOrPng = file.type === 'image/jpeg' || file.type === 'image/png';
  if (!isJpgOrPng) {
    message.error('不支持上传这个格式的图片，推荐JPG或PNG');
  }
  //校验图片大小
  const isLt2M = file.size / 1024 / 1024 < 2;
  if (!isLt2M) {
    message.error('不能上传超过2MB的图片');
  }
  return isJpgOrPng && isLt2M;
};
</script>
<style scoped>
.picture-upload :deep(.ant-upload) {
  width: 100% !important;
  height: 100% !important;
  min-width: 152px;
  min-height: 152px;
}
.avatar-uploader > .ant-upload {
  width: 128px;
  height: 128px;
}

.picture-upload img {
  max-width: 100%;
  max-height: 480px;
}
.ant-upload-select-picture-card i {
  font-size: 32px;
  color: #999;
}

.ant-upload-select-picture-card .ant-upload-text {
  margin-top: 8px;
  color: #666;
}
</style>
