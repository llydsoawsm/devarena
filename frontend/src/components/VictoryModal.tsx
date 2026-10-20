'use client';

import React, { useEffect } from 'react';
import { useRouter } from 'next/navigation';
import { useAuthStore } from '../stores/authStore';
import { useMatchStore } from '../stores/matchStore';

export function VictoryModal() {
  const router = useRouter();
  const { user, loadUser } = useAuthStore();
  const { matchResult, setMatchResult, reset } = useMatchStore();

  useEffect(() => {
    if (matchResult && user) {
      loadUser();
    }
  }, [matchResult, user, loadUser]);

  if (!matchResult) return null;

  const isWinner = matchResult.winner_id === user?.id;
  const isDraw = matchResult.is_draw;

  const handleReturnHome = () => {
    setMatchResult(null);
    reset();
    router.push('/');
  };

  return (
    <div className="fixed inset-0 z-50 flex items-center justify-center p-4 bg-black/70 backdrop-blur-xs">
      <div className="w-full max-w-sm p-6 rounded-lg bg-[#121215] border border-[#27272a] shadow-xl text-center space-y-5">
        <div>
          <span className="text-[10px] uppercase font-mono tracking-widest text-zinc-500 block mb-1">
            Match Concluded
          </span>
          <h2 className="text-2xl font-bold tracking-tight text-white">
            {isWinner ? 'Victory' : isDraw ? 'Draw' : 'Defeat'}
          </h2>
        </div>

        {/* Results Breakdown */}
        <div className="p-3.5 rounded bg-zinc-900 border border-zinc-800 space-y-2.5">
          {matchResult.participants?.map((p: any, idx: number) => {
            const isMe = p.user_id === user?.id;
            return (
              <div key={idx} className="flex items-center justify-between text-xs font-mono">
                <span className={isMe ? 'text-white font-medium' : 'text-zinc-500'}>
                  {p.name} {isMe && '(You)'}
                </span>
                <div className="flex items-center gap-1.5">
                  <span className="text-zinc-300 font-semibold">{p.rating_after} ELO</span>
                  <span
                    className={
                      p.rating_delta >= 0 ? 'text-emerald-400 font-bold' : 'text-rose-400 font-bold'
                    }
                  >
                    ({p.rating_delta >= 0 ? `+${p.rating_delta}` : p.rating_delta})
                  </span>
                </div>
              </div>
            );
          })}
        </div>

        <button
          onClick={handleReturnHome}
          className="w-full py-2 px-3 rounded bg-white hover:bg-zinc-200 text-black font-medium text-xs transition-colors"
        >
          Return to Arena
        </button>
      </div>
    </div>
  );
}
