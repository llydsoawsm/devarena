import { create } from 'zustand';
import { Match, OpponentState } from '../lib/types';

interface MatchResultData {
  winner_id?: number;
  is_draw: boolean;
  participants: Array<{
    user_id: number;
    name: string;
    result: string;
    rating_after: number;
    rating_delta: number;
  }>;
}

interface MatchState {
  isQueued: boolean;
  queuedMode: 'ranked' | 'casual' | 'tournament' | null;
  queueSeconds: number;
  activeMatch: Match | null;
  countdown: number;
  isMatchLive: boolean;
  opponentProgress: OpponentState | null;
  matchResult: MatchResultData | null;

  setQueued: (queued: boolean, mode?: 'ranked' | 'casual' | 'tournament') => void;
  incrementQueueTimer: () => void;
  setActiveMatch: (match: Match | null) => void;
  setCountdown: (countdown: number) => void;
  setMatchLive: (live: boolean) => void;
  setOpponentProgress: (progress: OpponentState | null) => void;
  setMatchResult: (result: MatchResultData | null) => void;
  reset: () => void;
}

export const useMatchStore = create<MatchState>((set) => ({
  isQueued: false,
  queuedMode: null,
  queueSeconds: 0,
  activeMatch: null,
  countdown: 10,
  isMatchLive: false,
  opponentProgress: null,
  matchResult: null,

  setQueued: (isQueued, mode = 'ranked') =>
    set({
      isQueued,
      queuedMode: isQueued ? mode : null,
      queueSeconds: 0,
    }),

  incrementQueueTimer: () =>
    set((state) => ({ queueSeconds: state.queueSeconds + 1 })),

  setActiveMatch: (activeMatch) =>
    set({
      activeMatch,
      isMatchLive: activeMatch?.status === 'live',
    }),

  setCountdown: (countdown) => set({ countdown }),

  setMatchLive: (isMatchLive) => set({ isMatchLive }),

  setOpponentProgress: (opponentProgress) => set({ opponentProgress }),

  setMatchResult: (matchResult) => set({ matchResult }),

  reset: () =>
    set({
      isQueued: false,
      queuedMode: null,
      queueSeconds: 0,
      activeMatch: null,
      countdown: 10,
      isMatchLive: false,
      opponentProgress: null,
      matchResult: null,
    }),
}));
