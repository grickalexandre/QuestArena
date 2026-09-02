import { useEffect, useState } from 'react'

type Send = (type: string, data?: unknown) => void

function isInspectShortcut(e: KeyboardEvent) {
  const key = e.key.toLowerCase()
  const ctrl = e.ctrlKey || e.metaKey
  if (e.key === 'F12') return true
  if (ctrl && e.shiftKey && (key === 'i' || key === 'j' || key === 'c' || key === 'k')) return true
  if (ctrl && key === 'u') return true
  return false
}

function desktopDevtoolsOpen() {
  if (window.matchMedia('(pointer: coarse)').matches) return false
  if (window.innerWidth < 768) return false
  const widthGap = window.outerWidth - window.innerWidth
  const heightGap = window.outerHeight - window.innerHeight
  return widthGap > 220 || heightGap > 220
}

/**
 * Blocks common inspect shortcuts on the play page and reports DevTools
 * to the teacher. This is a deterrent — scoring stays server-side.
 */
export function usePlayAntiInspect(send: Send, active: boolean) {
  const [inspecting, setInspecting] = useState(false)

  useEffect(() => {
    if (!active) {
      setInspecting(false)
      return
    }

    let last: boolean | null = null
    const report = (open: boolean) => {
      if (last === open) return
      last = open
      setInspecting(open)
      send('presence', { inspect: open })
    }

    const onKey = (e: KeyboardEvent) => {
      if (!isInspectShortcut(e)) return
      e.preventDefault()
      e.stopPropagation()
      report(true)
    }

    const onContext = (e: MouseEvent) => {
      e.preventDefault()
    }

    const detect = () => {
      report(desktopDevtoolsOpen())
    }

    document.addEventListener('keydown', onKey, true)
    document.addEventListener('contextmenu', onContext)
    window.addEventListener('resize', detect)
    const id = window.setInterval(detect, 1500)
    detect()

    return () => {
      document.removeEventListener('keydown', onKey, true)
      document.removeEventListener('contextmenu', onContext)
      window.removeEventListener('resize', detect)
      window.clearInterval(id)
      if (last === true) send('presence', { inspect: false })
    }
  }, [active, send])

  return inspecting
}
