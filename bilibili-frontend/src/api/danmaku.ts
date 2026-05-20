import request from './request'

/** 弹幕数据结构 */
export interface DanmakuItem {
  id: number
  video_id: number
  user_id: number
  content: string
  time_point: number
  color: string
  type: number
  created_at: string
}

/** 发送弹幕请求 */
export interface SendDanmakuPayload {
  content: string
  time_point: number
  color?: string
  type?: number
}

/** 获取某时间段弹幕列表 */
export async function getDanmaku(
  videoID: number,
  start: number,
  end: number,
): Promise<{ list: DanmakuItem[] }> {
  return request.get(`/videos/${videoID}/danmaku?start=${start}&end=${end}`)
}

/** 发送弹幕 */
export async function sendDanmaku(
  videoID: number,
  payload: SendDanmakuPayload,
): Promise<{ danmaku: DanmakuItem }> {
  return request.post(`/videos/${videoID}/danmaku`, payload)
}
