// =========================
// 首页模块 TypeScript 类型定义
// =========================

/** 分类导航项 */
export interface CategoryItem {
  id: string
  label: string
  code: string
}

/** Banner 数据 */
export interface BannerItem {
  id: number
  title: string
  subtitle: string
  gradientColors: string[] // 渐变色数组
}

/** UP 主信息 */
export interface UploaderInfo {
  name: string
  avatar?: string       // 真实头像 URL
  avatarColor: string   // 无头像时的 fallback 颜色
}

/** 视频数据 */
export interface VideoItem {
  id: number
  title: string
  cover_url?: string      // 真实封面图片 URL
  coverGradient: string[] // 封面渐变色（fallback）
  duration: string        // 如 "14:02"
  uploader: UploaderInfo
  views: number
  danmaku: number
  date: string            // 相对时间，如 "> 3h ago"
  category: string
  tags?: string[]         // 如 ['4K', 'HDR']
}

/** 热搜词 */
export interface HotSearchItem {
  id: number
  text: string
  heat: 'hot' | 'new' | 'normal'
}

/** 侧边栏菜单项 */
export interface MenuItem {
  id: string
  label: string
  icon: string            // SVG path d
  active?: boolean
  expandable?: boolean
  children?: MenuItem[]
}

/** 排行榜项 */
export interface RankItem {
  id: number
  title: string
  views: number
  rank: 1 | 2 | 3
}

/** 分区区块配置 */
export interface SectionConfig {
  title: string
  videos: VideoItem[]
  ranks: RankItem[]
}
