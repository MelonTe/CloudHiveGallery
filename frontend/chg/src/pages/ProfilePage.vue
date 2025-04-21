<template>
  <div id="profilePage" class="glass-effect">
    <div class="profile-container">
      <!-- 背景装饰 -->
      <div class="background-decoration">
        <div class="floating-note" v-for="(note, index) in notes" :key="index" :style="note.style">
          <span>{{ note.text }}</span>
        </div>
      </div>

      <!-- 用户信息卡片 -->
      <div class="profile-card">
        <div class="avatar-section">
          <a-upload name="avatar" :show-upload-list="false" :custom-request="handleAvatarUpload"
            :before-upload="beforeAvatarUpload">
            <a-avatar :size="120" :src="userInfo.userAvatar" class="profile-avatar" />
          </a-upload>
          <div class="user-basic-info">
            <h2>{{ userInfo.userName }}</h2>
            <p class="user-id">ID: {{ userInfo.id }}</p>
          </div>
        </div>

        <div class="profile-content">
          <a-form :model="editForm" :rules="rules" ref="formRef" layout="vertical" class="edit-form">
            <a-form-item label="昵称" name="userName">
              <a-input v-model:value="editForm.userName" placeholder="请输入昵称" :disabled="!isEditing" />
            </a-form-item>
            <a-form-item label="个人简介" name="userProfile">
              <a-textarea v-model:value="editForm.userProfile" placeholder="请输入个人简介" :rows="4" :disabled="!isEditing" />
            </a-form-item>
            <div class="form-actions">
              <a-button v-if="!isEditing" type="primary" @click="startEditing" class="edit-btn">
                编辑信息
              </a-button>
              <template v-else>
                <a-button type="primary" @click="saveProfile" class="save-btn">
                  保存
                </a-button>
                <a-button @click="cancelEditing" class="cancel-btn">
                  取消
                </a-button>
              </template>
            </div>
          </a-form>
        </div>
      </div>
    </div>
  </div>
</template>

<script setup lang="ts">
import { ref, reactive, onMounted, computed } from 'vue'
import { message } from 'ant-design-vue'
import { useLoginUserStore } from '../stores/useLoginUserStore'
import { postUserEdit, postUserAvatar } from '../api/user'
import type { FormInstance } from 'ant-design-vue'

interface ApiResponse<T = any> {
  code: number
  data: T
  msg: string
}

const loginUserStore = useLoginUserStore()
const formRef = ref<FormInstance>()
const isEditing = ref(false)
const originalForm = ref({})

// 用户信息
const userInfo = computed(() => loginUserStore.loginUser)

// 编辑表单
const editForm = reactive({
  userName: userInfo.value.userName,
  userProfile: userInfo.value.userProfile,
})

// 表单验证规则
const rules = {
  userName: [
    { required: true, message: '请输入昵称', trigger: 'blur' },
    { min: 2, max: 20, message: '昵称长度在2-20个字符之间', trigger: 'blur' },
  ],
  userProfile: [
    { max: 200, message: '个人简介不能超过200个字符', trigger: 'blur' },
  ],
}

// 开始编辑
const startEditing = () => {
  isEditing.value = true
  originalForm.value = { ...editForm }
}

// 取消编辑
const cancelEditing = () => {
  isEditing.value = false
  Object.assign(editForm, originalForm.value)
}

// 保存个人信息
const saveProfile = async () => {
  try {
    await formRef.value?.validate()
    const res = await postUserEdit({
      id: userInfo.value.id,
      userName: editForm.userName,
      userProfile: editForm.userProfile
    })
    const response = res.data as unknown as ApiResponse
    if (response.code === 0) {
      message.success('保存成功')
      loginUserStore.setLoginUser({
        ...loginUserStore.loginUser,
        ...editForm,
      })
      isEditing.value = false
    } else {
      message.error('保存失败：' + response.msg)
    }
  } catch (error) {
    console.error('保存失败：', error)
  }
}

// 上传头像
const handleAvatarUpload = async ({ file }: any) => {
  try {
    const res = await postUserAvatar({}, file)
    const response = res.data as unknown as ApiResponse<string>
    if (response.code === 0) {
      message.success('头像上传成功')
      loginUserStore.setLoginUser({
        ...loginUserStore.loginUser,
        userAvatar: response.data,
      })
      setTimeout(() => {
        window.location.reload()
      }, 800)
    } else {
      message.error('头像上传失败：' + response.msg)
    }
  } catch (error) {
    message.error('头像上传失败')
  }
}

// 头像上传前的验证
const beforeAvatarUpload = (file: File) => {
  const isJpgOrPng = file.type === 'image/jpeg' || file.type === 'image/png'
  if (!isJpgOrPng) {
    message.error('只能上传JPG/PNG格式的图片!')
  }
  const isLt2M = file.size / 1024 / 1024 < 2
  if (!isLt2M) {
    message.error('图片大小不能超过2MB!')
  }
  return isJpgOrPng && isLt2M
}

// 浮动音符效果
const notes = ref([
  { text: '♪', style: { top: '10%', left: '5%', animationDelay: '0s' } },
  { text: '♫', style: { top: '20%', left: '15%', animationDelay: '1s' } },
  { text: '♬', style: { top: '30%', left: '25%', animationDelay: '2s' } },
  { text: '♩', style: { top: '40%', left: '35%', animationDelay: '3s' } },
  { text: '♪', style: { top: '50%', left: '45%', animationDelay: '4s' } },
  { text: '♫', style: { top: '60%', left: '55%', animationDelay: '5s' } },
  { text: '♬', style: { top: '70%', left: '65%', animationDelay: '6s' } },
  { text: '♩', style: { top: '80%', left: '75%', animationDelay: '7s' } },
])
</script>

<style scoped>
#profilePage {
  min-height: calc(100vh - 130px);
  position: relative;
  overflow: hidden;
}

.profile-container {
  max-width: 800px;
  margin: 0 auto;
  position: relative;
}

.background-decoration {
  position: absolute;
  top: 0;
  left: 0;
  right: 0;
  bottom: 0;
  pointer-events: none;
  z-index: 0;
}

.floating-note {
  position: absolute;
  font-size: 24px;
  color: rgba(100, 212, 135, 0.3);
  animation: float 6s infinite ease-in-out;
}

@keyframes float {
  0% {
    transform: translateY(0) rotate(0deg);
  }

  50% {
    transform: translateY(-20px) rotate(10deg);
  }

  100% {
    transform: translateY(0) rotate(0deg);
  }
}

.profile-card {
  background: rgba(255, 255, 255, 0.95);
  border-radius: 12px;
  padding: 24px;
  box-shadow: 0 4px 12px rgba(0, 0, 0, 0.05);
  position: relative;
  z-index: 1;
}

.avatar-section {
  display: flex;
  align-items: center;
  margin-bottom: 24px;
}

.profile-avatar {
  border: 3px solid #64d487;
  box-shadow: 0 4px 12px rgba(100, 212, 135, 0.2);
  cursor: pointer;
  transition: all 0.3s;
}

.profile-avatar:hover {
  transform: scale(1.05);
  box-shadow: 0 6px 16px rgba(100, 212, 135, 0.3);
}

.user-basic-info {
  margin-left: 24px;
}

.user-basic-info h2 {
  margin: 0;
  color: #333;
  font-size: 24px;
}

.user-id {
  color: #999;
  margin: 8px 0 0;
}

.profile-content {
  margin-top: 24px;
}

.edit-form {
  max-width: 500px;
}

.form-actions {
  margin-top: 24px;
  display: flex;
  gap: 12px;
}

.edit-btn,
.save-btn,
.cancel-btn {
  border-radius: 6px;
  padding: 0 24px;
  height: 40px;
}

.edit-btn {
  background: #64d487;
  border-color: #64d487;
}

.edit-btn:hover {
  background: #43b16a;
  border-color: #43b16a;
}

.save-btn {
  background: #64d487;
  border-color: #64d487;
}

.save-btn:hover {
  background: #43b16a;
  border-color: #43b16a;
}

.cancel-btn {
  border-color: #d9d9d9;
  color: #666;
}

.cancel-btn:hover {
  border-color: #64d487;
  color: #64d487;
}

/* 移除响应式边距设置，由 BasicLayout 控制 */
@media (max-width: 992px) {
  #profilePage {
    margin-left: 0;
  }
}

@media (max-width: 576px) {
  #profilePage {
    margin-left: 0;
  }
}
</style>