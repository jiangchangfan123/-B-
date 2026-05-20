declare module 'dplayer' {
  interface DPlayerVideo {
    url: string
    pic?: string
    type?: string
    customType?: Record<string, (video: HTMLVideoElement, player: any) => void>
  }

  interface DPlayerOptions {
    container: HTMLElement | null
    video: DPlayerVideo
    theme?: string
    autoplay?: boolean
    lang?: string
    screenshot?: boolean
    hotkey?: boolean
    preload?: string
    volume?: number
    mutex?: boolean
    contextmenu?: Array<{ text: string; link?: string; click?: (player: any) => void }>
    highlight?: Array<{ time: number; text: string }>
  }

  class DPlayer {
    constructor(options: DPlayerOptions)
    video: HTMLVideoElement
    destroy(): void
    seek(time: number): void
    notice(text: string, time?: number, opacity?: number): void
    switchVideo(video: DPlayerVideo, danmaku?: any): void
    on(event: string, handler: (...args: any[]) => void): void
    off(event: string, handler: (...args: any[]) => void): void
  }

  export default DPlayer
}
