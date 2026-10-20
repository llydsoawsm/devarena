'use client';

import React, { useEffect } from 'react';
import { useMatchStore } from '../stores/matchStore';
import { matchmakingApi } from '../lib/api';
import { toast } from 'sonner';

export function MatchmakingModal() {
  const { isQueued, queuedMode, queueSeconds, setQueued, incrementQueueTimer } = useMatchStore();

  useEffect(() => {
    let interval: NodeJS.Timeout;
    if (isQueued) {
      interval = setInterval(() => {
        incrementQueueTimer();
      }, 1000);
    }
    return () => clearInterval(interval);
  }, [isQueued, incrementQueueTimer]);

  if (!isQueued) return null;

  const handleCancel = async () => {
    try {
      await matchmakingApi.dequeue();
      setQueued(false);
      toast.info('Matchmaking cancelled');
    } catch {
      setQueued(false);
    }
  };

  const formatTimer = (seconds: number) => {
    const mins = Math.floor(seconds / 60);
    const secs = seconds % 60;
    return `${mins.toString().padStart(2, '0')}:${secs.toString().padStart(2, '0')}`;
  };

  return (
    <div className="fixed inset-0 z-50 flex items-center justify-center p-4 bg-black/60 backdrop-blur-xs">
      <div className="w-full max-w-sm p-6 rounded-lg bg-[#121215] border border-[#27272a] shadow-xl text-center space-y-5">
        <div className="w-8 h-8 mx-auto border-2 border-zinc-700 border-t-white rounded-full animate-spin" />

        <div className="space-y-1">
          <h3 className="text-base font-semibold text-white">Searching for Opponent</h3>
          <p className="text-xs text-zinc-400 capitalize">
            Mode: {queuedMode} 1v1
          </p>
        </div>

        <div className="p-3 rounded bg-zinc-900 border border-zinc-800 flex items-center justify-between text-xs font-mono">
          <span className="text-zinc-500">Wait Time</span>
          <span className="text-white font-medium">{formatTimer(queueSeconds)}</span>
        </div>

        <button
          onClick={handleCancel}
          className="w-full py-2 px-3 rounded bg-zinc-800 hover:bg-zinc-700 text-zinc-300 text-xs font-medium transition-colors border border-zinc-700"
        >
          Cancel Search
        </button>
      </div>
    </div>
  );
}
