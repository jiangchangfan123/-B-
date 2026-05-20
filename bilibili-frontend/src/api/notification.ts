import request from './request'

export interface NotificationItem {
  id: number
  type: number
  title: string
  content: string
  related_id: number
  is_read: boolean
  created_at: string
  trigger_user?: {
    id: number
    username: string
    nickname: string
    avatar: string
  }
}

export interface NotificationListResponse {
  list: NotificationItem[]
  total: number
  page: number
  size: number
}

export async function getNotifications(page = 1, size = 20, unreadOnly = false): Promise<NotificationListResponse> {
  const params = new URLSearchParams()
  params.append('page', String(page))
  params.append('size', String(size))
  if (unreadOnly) params.append('unread_only', 'true')
  const res = await request.get(`/users/me/notifications?${params.toString()}`)
  return res
}

export async function getUnreadCount(): Promise<number> {
  const res = await request.get('/users/me/notifications/unread-count')
  return res.count
}

export async function markAsRead(id: number): Promise<void> {
  await request.put(`/users/me/notifications/${id}/read`)
}

export async function markAllAsRead(): Promise<void> {
  await request.put('/users/me/notifications/read-all')
}

export async function deleteNotification(id: number): Promise<void> {
  await request.delete(`/users/me/notifications/${id}`)
}
