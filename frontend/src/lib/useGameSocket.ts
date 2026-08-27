import { useCallback, useEffect, useRef, useState } from 'react'
import { wsUrl } from './api'

export type WsMessage = { type: string; data?: unknown }

type Handler = (data: unknown) => void

export function useGameSocket() {
  const wsRef = useRef<WebSocket | null>(null)
  const handlers = useRef(new Map<string, Handler>())
  const openHandlers = useRef(new Set<() => void>())
  const closeHandlers = useRef(new Set<() => void>())
  const wantedRef = useRef(false)
  const retryRef = useRef(0)
  const retryTimer = useRef<ReturnType<typeof setTimeout> | null>(null)
  const connectRef = useRef<() => WebSocket>(() => {
    throw new Error('socket not ready')
  })
  const [connected, setConnected] = useState(false)
  const [lastMessage, setLastMessage] = useState<WsMessage | null>(null)

  const clearRetry = useCallback(() => {
    if (retryTimer.current) {
      clearTimeout(retryTimer.current)
      retryTimer.current = null
    }
  }, [])

  const attach = useCallback((ws: WebSocket) => {
    ws.onopen = () => {
      retryRef.current = 0
      setConnected(true)
      openHandlers.current.forEach((h) => h())
    }
    ws.onclose = () => {
      if (wsRef.current !== ws) return
      wsRef.current = null
      setConnected(false)
      closeHandlers.current.forEach((h) => h())
      if (!wantedRef.current) return
      const delay = Math.min(8000, 400 * 2 ** retryRef.current)
      retryRef.current = Math.min(retryRef.current + 1, 5)
      clearRetry()
      retryTimer.current = setTimeout(() => {
        if (wantedRef.current) connectRef.current()
      }, delay)
    }
    ws.onerror = () => {
      /* onclose follows */
    }
    ws.onmessage = (ev) => {
      try {
        const msg = JSON.parse(ev.data) as { type: string; data?: string | object }
        let data: unknown = msg.data
        if (typeof msg.data === 'string') {
          try {
            data = JSON.parse(msg.data)
          } catch {
            data = msg.data
          }
        }
        const parsed = { type: msg.type, data }
        setLastMessage(parsed)
        handlers.current.get(msg.type)?.(data)
        handlers.current.get('*')?.(parsed)
      } catch {
        /* ignore */
      }
    }
  }, [clearRetry])

  const connect = useCallback(() => {
    wantedRef.current = true
    const cur = wsRef.current
    if (cur && (cur.readyState === WebSocket.OPEN || cur.readyState === WebSocket.CONNECTING)) {
      return cur
    }
    const ws = new WebSocket(wsUrl())
    wsRef.current = ws
    attach(ws)
    return ws
  }, [attach])
  connectRef.current = connect

  const disconnect = useCallback(() => {
    wantedRef.current = false
    clearRetry()
    const ws = wsRef.current
    wsRef.current = null
    setConnected(false)
    ws?.close()
  }, [clearRetry])

  const send = useCallback(
    (type: string, data?: unknown) => {
      const ws = connect()
      const payload = JSON.stringify({ type, data })
      if (ws.readyState === WebSocket.OPEN) {
        ws.send(payload)
        return
      }
      const flush = () => {
        if (ws.readyState === WebSocket.OPEN) ws.send(payload)
      }
      ws.addEventListener('open', flush, { once: true })
    },
    [connect],
  )

  const on = useCallback((type: string, handler: Handler) => {
    handlers.current.set(type, handler)
    return () => {
      if (handlers.current.get(type) === handler) {
        handlers.current.delete(type)
      }
    }
  }, [])

  const onOpen = useCallback((handler: () => void) => {
    openHandlers.current.add(handler)
    if (wsRef.current?.readyState === WebSocket.OPEN) handler()
    return () => {
      openHandlers.current.delete(handler)
    }
  }, [])

  const onClose = useCallback((handler: () => void) => {
    closeHandlers.current.add(handler)
    return () => {
      closeHandlers.current.delete(handler)
    }
  }, [])

  useEffect(() => () => disconnect(), [disconnect])

  return { connected, lastMessage, connect, disconnect, send, on, onOpen, onClose }
}
