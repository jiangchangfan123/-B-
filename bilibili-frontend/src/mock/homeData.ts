import type {
  BannerItem,
  CategoryItem,
  HotSearchItem,
  MenuItem,
  RankItem,
  VideoItem,
} from '../types/home'

// =========================
// 模拟数据
// =========================

export const bannerList: BannerItem[] = [
  {
    id: 1,
    title: '深空探索计划 // 启航',
    subtitle: 'SYSTEM_UPDATE // v2.5.0',
    gradientColors: ['#1a0b2e', '#4d1b7b', '#00f0ff'],
  },
  {
    id: 2,
    title: '量子计算机入门指南',
    subtitle: 'TECH_BRIEF // Q3-2025',
    gradientColors: ['#0f1117', '#b829dd', '#4d6bfa'],
  },
  {
    id: 3,
    title: '极光城市 // 夜间模式',
    subtitle: 'VISUAL_EXPERIMENT // NEON',
    gradientColors: ['#08090d', '#00f0ff', '#b829dd'],
  },
  {
    id: 4,
    title: 'AI 生成艺术的未来',
    subtitle: 'CREATIVE_LAB // GEN-AI-07',
    gradientColors: ['#14161f', '#4d6bfa', '#00f0ff'],
  },
  {
    id: 5,
    title: '星际旅行日记 // 第七章',
    subtitle: 'LOG_ENTRY // SECTOR-7G',
    gradientColors: ['#0a0a0f', '#b829dd', '#4d6bfa'],
  },
]

export const categoryList: CategoryItem[] = [
  { id: 'cine', label: '影视', code: '[CINE]' },
  { id: 'game', label: '游戏', code: '[GAME]' },
  { id: 'acgn', label: '二次元', code: '[ACGN]' },
  { id: 'chef', label: '美食', code: '[CHEF]' },
  { id: 'docu', label: '纪录片', code: '[DOCU]' },
  { id: 'show', label: '综艺', code: '[SHOW]' },
]

export const hotSearchList: HotSearchItem[] = [
  { id: 1, text: '量子物理学新突破', heat: 'hot' },
  { id: 2, text: '人工智能诊断系统', heat: 'hot' },
  { id: 3, text: '火星移民计划启动', heat: 'new' },
  { id: 4, text: '虚拟现实头显测评', heat: 'normal' },
  { id: 5, text: '光学隐身技术发展', heat: 'normal' },
  { id: 6, text: '机械义肢控制系统', heat: 'normal' },
  { id: 7, text: '深海生物发光科学', heat: 'normal' },
  { id: 8, text: '基因编辑基础入门', heat: 'new' },
  { id: 9, text: '太阳风暴监测报告', heat: 'normal' },
  { id: 10, text: '银河中心黑洞证据', heat: 'normal' },
]

const avatarColors = ['#00f0ff', '#b829dd', '#4d6bfa', '#e4e5eb', '#8b8fa3']
const coverPalettes = [
  ['#14161f', '#b829dd'],
  ['#14161f', '#00f0ff'],
  ['#14161f', '#4d6bfa'],
  ['#0f1117', '#b829dd'],
  ['#0f1117', '#00f0ff'],
  ['#0a0a0f', '#4d6bfa'],
]

function randomPalette(): string[] {
  return coverPalettes[Math.floor(Math.random() * coverPalettes.length)]
}

function randomAvatarColor(): string {
  return avatarColors[Math.floor(Math.random() * avatarColors.length)]
}

function formatDuration(seconds: number): string {
  const m = Math.floor(seconds / 60)
    .toString()
    .padStart(2, '0')
  const s = (seconds % 60).toString().padStart(2, '0')
  return `${m}:${s}`
}

const videoTitles = [
  '[4K] 极光城市夜间模式 // 游戏开发实录',
  '量子计算机原理入门 // 第一课',
  '[HDR] 星际旅行模拟器 // 新手教程',
  '深海探索器传回的影像',
  'AI 图像生成的后座力 // 艺术论坛',
  '[4K] 机械义肢实战演示',
  '火星基地建设进度报告',
  '光学隐身实验 // 全过程记录',
  '银河系全景模拟 // 60fps',
  '基因编辑技术最新进展',
  '[HDR] 幻想世界构建 // 场景设计',
  '太阳风暴监测数据分析',
  'VR 深海体验 // 沉浸式一周',
  '城市交通模拟系统 // 开源项目',
  '无人机跨大洋飞行 // 实时监控',
  '深空探测器 // 系统测试',
  '[4K] 虚拟现实影像质量测试',
  '星际通信技术解密',
]

const uploaderNames = [
  'NEBULA_LABS',
  'QUANTUM_DEV',
  'VOID_WALKER',
  'DEEP_DIVE',
  'SYNTH_MIND',
  'CYBER_CRAFT',
  'MARS_BASE',
  'OPTICS_LAB',
  'GALAXY_SIM',
  'GENE_EDIT',
]

export const videoList: VideoItem[] = videoTitles.map((title, i) => {
  const views = Math.floor(Math.random() * 500000) + 1000
  const danmaku = Math.floor(views * 0.02)
  return {
    id: i + 1,
    title,
    coverGradient: randomPalette(),
    duration: formatDuration(Math.floor(Math.random() * 1800) + 120),
    uploader: {
      name: uploaderNames[i % uploaderNames.length],
      avatarColor: randomAvatarColor(),
    },
    views,
    danmaku,
    date: `> ${Math.floor(Math.random() * 23 + 1)}h ago`,
    category: categoryList[i % categoryList.length].id,
    tags: i % 3 === 0 ? ['4K'] : i % 5 === 0 ? ['HDR'] : undefined,
  }
})

export const menuList: MenuItem[] = [
  {
    id: 'home',
    label: '首页',
    icon: 'M3 9l9-7 9 7v11a2 2 0 01-2 2H5a2 2 0 01-2-2V9z',
    active: true,
  },
  {
    id: 'favorites',
    label: '收藏',
    icon: 'M5 4a2 2 0 012-2h6a2 2 0 012 2v14l-5-2.5L5 18V4z',
    active: false,
  },
]

export const systemMenuList: MenuItem[] = [
  {
    id: 'notifications',
    label: '消息',
    icon: 'M15 17h5l-1.405-1.405A2.032 2.032 0 0118 14.158V11a6.002 6.002 0 00-4-5.659V5a2 2 0 10-4 0v.341C7.67 6.165 6 8.388 6 11v3.159c0 .538-.214 1.055-.595 1.436L4 17h5m6 0v1a3 3 0 11-6 0v-1m6 0H9',
    active: false,
  },
  {
    id: 'personal',
    label: '个人中心',
    icon: 'M20 21v-2a4 4 0 00-4-4H8a4 4 0 00-4 4v2M12 11a4 4 0 100-8 4 4 0 000 8z',
    active: false,
  },
]

export function formatNumber(n: number): string {
  if (n >= 10000) {
    return (n / 10000).toFixed(1) + 'W'
  }
  if (n >= 1000) {
    return (n / 1000).toFixed(1) + 'K'
  }
  return String(n)
}

export const allSectionVideos = videoList

export const gameSectionVideos = videoList.filter((v) => v.category === 'game')

export function generateRanks(videos: VideoItem[]): RankItem[] {
  const sorted = [...videos].sort((a, b) => b.views - a.views).slice(0, 3)
  return sorted.map((v, i) => ({
    id: v.id,
    title: v.title,
    views: v.views,
    rank: (i + 1) as 1 | 2 | 3,
  }))
}
