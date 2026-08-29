import { useEffect, useState } from 'react'

type Send = (type: string, data?: unknown) => void

/**
 * Reports when the student leaves the quiz tab/app.
 * Relies on Page Visibility + pagehide/pageshow (iOS) and freeze/resume
 * (Android). Does not use window.blur — the mobile keyboard would false-trigger.
 */
export function usePlayPresence(send: Send, active: boolean) {
  const [away, setAway] = useState(false)

  useEffect(() => {
    if (!active) {
      setAway(false)
      return
    }

    let lastHidden: boolean | null = null

    const report = (hidden: boolean) => {
      if (lastHidden === hidden) return
      lastHidden = hidden
      setAway(hidden)
      send('presence', { hidden })
    }

    const syncFromVisibility = () => {
      report(document.visibilityState !== 'visible')
    }

    const onPageHide = () => report(true)
    const onPageShow = () => report(false)
    const onFreeze = () => report(true)
    const onResume = () => report(false)

    document.addEventListener('visibilitychange', syncFromVisibility)
    window.addEventListener('pagehide', onPageHide)
    window.addEventListener('pageshow', onPageShow)
    document.addEventListener('freeze', onFreeze)
    document.addEventListener('resume', onResume)

    syncFromVisibility()

    return () => {
      document.removeEventListener('visibilitychange', syncFromVisibility)
      window.removeEventListener('pagehide', onPageHide)
      window.removeEventListener('pageshow', onPageShow)
      document.removeEventListener('freeze', onFreeze)
      document.removeEventListener('resume', onResume)
    }
  }, [active, send])

  return away
}
