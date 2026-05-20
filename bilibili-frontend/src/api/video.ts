import request from './request'
import type {
  UploadResponse,
  TranscodeStatusResponse,
  VideoListResponse,
  VideoDetail,
} from '../types/video'

/** 更新视频 */
export async function updateVideo(id: number, formData: FormData): Promise<void> {
  return request.put(`/videos/${id}`, formData, {
    headers: {
      'Content-Type': 'multipart/form-data',
    },
  })
}
const USE_MOCK = false

let mockPollCount = 0

/**视频上传 */
export async function uploadVideo(formData: FormData): Promise<UploadResponse> {
  if (USE_MOCK) {
    await delay(2000)
    return {
      id: 10001,
      title: formData.get('title') as string || '未命名视频',
      cover_url: '',
      status: 2,
      transcode_status: 1,
      created_at: new Date().toISOString(),
    }
  }
  return request.post('/videos', formData, {
    headers: {
      'Content-Type': 'multipart/form-data',
    },
  })
}

/**查询转码状态 */
export async function getTranscodeStatus(id: number): Promise<TranscodeStatusResponse> {
  if (USE_MOCK) {
    await delay(2000)
    mockPollCount++
    if (mockPollCount < 4) {
      return { transcode_status: 1, transcoded_url: '' }
    }
    return {
      transcode_status: 2,
      transcoded_url: 'https://www.w3schools.com/html/mov_bbb.mp4',
    }
  }
  return request.get(`/videos/${id}/transcode`)
}

/**视频列表 */
export async function getVideoList(
  page = 1,
  size = 20,
  category?: string,
  sort = 'new',
): Promise<VideoListResponse> {
  if (USE_MOCK) {
    await delay(300)
    return {
      list: generateMockVideoList(category),
      total: 8,
      page,
      size,
    }
  }
  let url = `/videos?page=${page}&size=${size}&sort=${sort}`
  if (category && category !== 'all') {
    url += `&category=${category}`
  }
  return request.get(url)
}

/**视频详情 */
export async function getVideoDetail(id: number): Promise<VideoDetail> {
  if (USE_MOCK) {
    await delay(300)
    return generateMockVideoDetail(id)
  }
  return request.get(`/videos/${id}`)
}

/**删除视频 */
export async function deleteVideo(id: number): Promise<void> {
  if (USE_MOCK) {
    await delay(500)
    return
  }
  return request.delete(`/videos/${id}`)
}

/** 点赞/取消点赞 */
export async function likeVideo(id: number): Promise<{ liked: boolean; count: number }> {
  return request.post(`/videos/${id}/like`)
}

/** 查询点赞状态 */
export async function getLikeStatus(id: number): Promise<{ liked: boolean; count: number }> {
  return request.get(`/videos/${id}/like/status`)
}

/** 收藏/取消收藏 */
export async function favoriteVideo(id: number): Promise<{ favorited: boolean }> {
  return request.post(`/videos/${id}/favorite`)
}

/** 查询收藏状态 */
export async function getFavoriteStatus(id: number): Promise<{ favorited: boolean }> {
  return request.get(`/videos/${id}/favorite/status`)
}

/** 我的收藏列表 */
export async function getFavorites(page = 1, size = 20): Promise<VideoListResponse> {
  return request.get(`/users/me/favorites?page=${page}&size=${size}`)
}
function delay(ms: number) {
  return new Promise(resolve => setTimeout(resolve, ms))
}

function generateMockVideoList(category?: string): any[] {
  const all = [
    { id: 1, title: '深空探索者：未知星球的秘密', cover_url: '', category: 'scifi', view_count: 12500, like_count: 342, created_at: '2026-05-18T10:00:00Z', user: { id: 1, username: 'nebula', nickname: '深空漫步者', avatar: '' } },
    { id: 2, title: 'FFmpeg 实战：从入门到放弃', cover_url: '', category: 'tech', view_count: 8900, like_count: 256, created_at: '2026-05-17T14:30:00Z', user: { id: 2, username: 'coder', nickname: '码农日记', avatar: '' } },
    { id: 3, title: '赛博朋克 2077 最新 MOD 测试', cover_url: '', category: 'game', view_count: 23400, like_count: 891, created_at: '2026-05-16T09:00:00Z', user: { id: 3, username: 'gamer', nickname: '夜之城游侠', avatar: '' } },
    { id: 4, title: '合成波 电子音乐制作教程', cover_url: '', category: 'music', view_count: 5600, like_count: 178, created_at: '2026-05-15T20:00:00Z', user: { id: 4, username: 'beat', nickname: '音波制作人', avatar: '' } },
    { id: 5, title: '末日废土生存指南', cover_url: '', category: 'life', view_count: 4200, like_count: 134, created_at: '2026-05-14T16:00:00Z', user: { id: 5, username: 'survivor', nickname: '废土行者', avatar: '' } },
    { id: 6, title: '人工智能会取代程序员吗？', cover_url: '', category: 'tech', view_count: 18900, like_count: 567, created_at: '2026-05-13T11:00:00Z', user: { id: 2, username: 'coder', nickname: '码农日记', avatar: '' } },
  ]
  if (category && category !== 'all') {
    return all.filter(v => v.category === category)
  }
  return all
}

function generateMockVideoDetail(id: number): VideoDetail {
  return {
    id,
    title: '深空探索者：未知星球的秘密',
    description: '这是一个关于深空探索的视频，展示了未知星球的神秘景象。\n\n视频中包含了大量的 CGI 特效和科幻元素，希望大家喜欢。',
    cover_url: '',
    video_url: 'https://example.com/raw.mp4',
    transcoded_url: 'https://www.w3schools.com/html/mov_bbb.mp4',
    category: 'scifi',
    status: 1,
    transcode_status: 2,
    view_count: 12500,
    like_count: 342,
    comment_count: 56,
    danmaku_count: 128,
    created_at: '2026-05-18T10:00:00Z',
    user_info: {
      id: 1,
      username: 'nebula',
      nickname: '深空漫步者',
      avatar: '',
    },
  }
}
