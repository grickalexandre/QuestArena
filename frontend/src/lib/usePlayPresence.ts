import { useEffect, useState } from 'react'

type Send = (type: string, data?: unknown) => void

/** Screen lock / idle typically happens after this much without a tap. */
const IDLE_MS = 20_000

/**
 * Reports when the student leaves the quiz tab/app.
 * Screen lock (phone inactivity) is sent as idle and must not cost the question.
 * Relies on Page Visibility + pagehide/pageshow (iOS) and freeze/resume
 * (Android). Does not use window.blur — the mobile keyboard would false-trigger.
 */
export function usePlayPresence(send: Send, active: boolean, roundKey = '') {
  const [away, setAway] = useState(false)
  const [idle, setIdle] = useState(false)

  useEffect(() => {
    if (!active) {
      setAway(false)
      setIdle(false)
      return
    }

    let lastHidden: boolean | null = null
    let lastIdle = false
    let lastInput = Date.now()

    const bump = () => {
      lastInput = Date.now()
    }
    bump()

    const report = (hidden: boolean) => {
      const idleHide = hidden && Date.now() - lastInput >= IDLE_MS
      if (lastHidden === hidden && lastIdle === idleHide) return
      lastHidden = hidden
      lastIdle = idleHide
      setAway(hidden)
      setIdle(idleHide)
      send('presence', { hidden, idle: idleHide })
    }

    const syncFromVisibility = () => {
      report(document.visibilityState !== 'visible')
    }

    const onPageHide = () => report(true)
    const onPageShow = () => {
      bump()
      report(false)
    }
    const onFreeze = () => report(true)
    const onResume = () => {
      bump()
      report(false)
    }

    document.addEventListener('visibilitychange', syncFromVisibility)
    window.addEventListener('pagehide', onPageHide)
    window.addEventListener('pageshow', onPageShow)
    document.addEventListener('freeze', onFreeze)
    document.addEventListener('resume', onResume)
    window.addEventListener('pointerdown', bump, true)
    window.addEventListener('keydown', bump, true)
    window.addEventListener('touchstart', bump, { capture: true, passive: true })

    syncFromVisibility()

    return () => {
      document.removeEventListener('visibilitychange', syncFromVisibility)
      window.removeEventListener('pagehide', onPageHide)
      window.removeEventListener('pageshow', onPageShow)
      document.removeEventListener('freeze', onFreeze)
      document.removeEventListener('resume', onResume)
      window.removeEventListener('pointerdown', bump, true)
      window.removeEventListener('keydown', bump, true)
      window.removeEventListener('touchstart', bump, true)
    }
  }, [active, send, roundKey])

  return { away, idle }
}
