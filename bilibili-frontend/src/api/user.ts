import request from './request'
import type { UserProfile, UpdateProfileForm, UpdatePasswordForm, MyVideoItem, HistoryItem } from '../types/auth'

/**获取当前用户详情 */
export async function getMeDetail(): Promise<UserProfile> {
  return request.get('/users/me')
}

/**更新个人资料 */
export async function updateProfile(data: UpdateProfileForm): Promise<void> {
  return request.put('/users/me', data)
}

/**修改密码 */
export async function updatePassword(data: UpdatePasswordForm): Promise<void> {
  return request.put('/users/me/password', {
    oldPassword: data.oldPassword,
    newPassword: data.newPassword,
  })
}

/**上传头像 */
export async function uploadAvatar(file: File): Promise<{ avatar: string }> {
  const formData = new FormData()
  formData.append('file', file)
  return request.post('/users/me/avatar', formData, {
    headers: {
      'Content-Type': 'multipart/form-data',
    },
  })
}

/**获取我的视频 */
export async function getMyVideos(page = 1, size = 20): Promise<{ list: MyVideoItem[]; total: number; page: number; size: number }> {
  return request.get(`/users/me/videos?page=${page}&size=${size}`)
}

/**获取播放历史 */
export async function getHistory(page = 1, size = 20): Promise<{ list: HistoryItem[]; total: number; page: number; size: number }> {
  return request.get(`/users/me/history?page=${page}&size=${size}`)
}
