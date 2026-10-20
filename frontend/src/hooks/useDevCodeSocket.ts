'use client';

import { useEffect, useRef, useCallback } from 'react';
import { useRouter } from 'next/navigation';
import { useAuthStore } from '../stores/authStore';
import { useMatchStore } from '../stores/matchStore';
import { toast } from 'sonner';

export function useDevCodeSocket(matchId?: string) {
  const router = useRouter();
  const socketRef = useRef<WebSocket | null>(null);
  const { token, user } = useAuthStore();
  const {
    setCountdown,
    setMatchLive,
    setOpponentProgress,
    setMatchResult,
    setQueued,
  } = useMatchStore();

  const connect = useCallback(() => {
    if (typeof window === 'undefined') return;

    const baseWs = process.env.NEXT_PUBLIC_WS_URL || 'ws://localhost:8080/ws';
    const params = new URLSearchParams();
    if (token) params.set('token', token);
    if (matchId) params.set('match_id', matchId);

    const wsUrl = `${baseWs}?${params.toString()}`;
    const ws = new WebSocket(wsUrl);
    socketRef.current = ws;

    ws.onopen = () => {
      // If room specified, join room
      if (matchId) {
        ws.send(JSON.stringify({ event: 'join_room', payload: matchId }));
      }
    };

    ws.onmessage = (event) => {
      try {
        const data = JSON.parse(event.data);
        const { event: evtName, payload } = data;

        switch (evtName) {
          case 'match.created':
            setQueued(false);
            toast.success('⚔️ Match Found! Entering Arena...', { duration: 3000 });
            if (payload?.match_id) {
              router.push(`/match/${payload.match_id}`);
            }
            break;

          case 'match.countdown_tick':
            if (typeof payload?.countdown === 'number') {
              setCountdown(payload.countdown);
            }
            break;

          case 'match.started':
            setCountdown(0);
            setMatchLive(true);
            toast.info('🔥 Match is LIVE! Code fast!', { duration: 2500 });
            break;

          case 'participant.state':
            if (payload?.user_id && payload.user_id !== user?.id) {
              setOpponentProgress(payload);
            }
            break;

          case 'match.submission_result':
            if (payload?.user_id !== user?.id) {
              toast.info(`⚡ Opponent tested solution: ${payload.passed_tests}/${payload.total_tests} passed`);
            }
            break;

          case 'match.finished':
            setMatchLive(false);
            setMatchResult(payload);
            break;
        }
      } catch (err) {
        console.error('WebSocket parse error:', err);
      }
    };

    ws.onclose = () => {
      // Reconnect after brief backoff if not intentionally closed
      setTimeout(() => {
        if (socketRef.current === ws) {
          connect();
        }
      }, 3000);
    };

    // Heartbeat keepalive
    const heartbeatInterval = setInterval(() => {
      if (ws.readyState === WebSocket.OPEN) {
        ws.send(JSON.stringify({ event: 'participant.heartbeat' }));
      }
    }, 15000);

    return () => {
      clearInterval(heartbeatInterval);
      ws.close();
    };
  }, [token, matchId, user?.id, router, setCountdown, setMatchLive, setOpponentProgress, setMatchResult, setQueued]);

  useEffect(() => {
    const cleanup = connect();
    return () => {
      if (cleanup) cleanup();
    };
  }, [connect]);

  const sendParticipantState = (state: { passed_tests: number; total_tests: number; verdict?: string }) => {
    if (socketRef.current?.readyState === WebSocket.OPEN && user) {
      socketRef.current.send(
        JSON.stringify({
          event: 'participant.state',
          payload: {
            user_id: user.id,
            ...state,
          },
        })
      );
    }
  };

  return {
    socket: socketRef.current,
    sendParticipantState,
  };
}
