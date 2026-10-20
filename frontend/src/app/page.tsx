'use client';

import React, { useEffect, useState } from 'react';
import { useAuthStore } from '../stores/authStore';
import { useMatchStore } from '../stores/matchStore';
import { matchmakingApi, problemsApi, matchesApi } from '../lib/api';
import { Problem, User } from '../lib/types';
import { useDevCodeSocket } from '../hooks/useDevCodeSocket';
import { MatchmakingModal } from '../components/MatchmakingModal';
import { AuthModal } from '../components/AuthModal';
import { Swords, Code2, Users, Trophy, ArrowRight, ChevronRight } from 'lucide-react';
import { toast } from 'sonner';

export default function ArenaLobby() {
  const { user, loadUser } = useAuthStore();
  const { setQueued, isQueued } = useMatchStore();

  const [problems, setProblems] = useState<Problem[]>([]);
  const [leaderboard, setLeaderboard] = useState<User[]>([]);
  const [authModalOpen, setAuthModalOpen] = useState(false);
  const [loadingAction, setLoadingAction] = useState<string | null>(null);

  useDevCodeSocket();

  useEffect(() => {
    loadUser();
    problemsApi.list().then((res) => setProblems(res.data.problems || []));
    matchesApi.leaderboard().then((res) => setLeaderboard(res.data.leaderboard || []));
  }, [loadUser]);

  const handleStartQueue = async (mode: 'ranked' | 'casual' | 'tournament') => {
    if (!user) {
      setAuthModalOpen(true);
      return;
    }

    setLoadingAction(mode);
    try {
      await matchmakingApi.enqueue(mode);
      setQueued(true, mode);
      toast.success(`Entered ${mode} matchmaking`);
    } catch (err: any) {
      toast.error(err.response?.data?.error || 'Failed to enter queue');
    } finally {
      setLoadingAction(null);
    }
  };

  const getDifficultyClass = (diff: string) => {
    switch (diff) {
      case 'easy':
        return 'text-emerald-400 bg-emerald-950/40 border-emerald-800/50';
      case 'medium':
        return 'text-amber-400 bg-amber-950/40 border-amber-800/50';
      case 'hard':
        return 'text-rose-400 bg-rose-950/40 border-rose-800/50';
      default:
        return 'text-zinc-400 bg-zinc-900 border-zinc-800';
    }
  };

  return (
    <div className="min-h-screen py-10 px-4 sm:px-6 max-w-6xl mx-auto space-y-12">
      {/* Minimal Header */}
      <div className="space-y-3 pb-8 border-b border-[#27272a]">
        <h1 className="text-3xl sm:text-4xl font-bold tracking-tight text-white">
          Competitive Coding Arena
        </h1>
        <p className="text-sm text-zinc-400 max-w-2xl leading-relaxed">
          Real-time 1v1 developer arena. Queue against players near your rating, solve algorithmic challenges in sandboxed runners, and climb the ELO ladder.
        </p>

        <div className="pt-3 flex items-center gap-3">
          <button
            onClick={() => handleStartQueue('ranked')}
            disabled={isQueued || loadingAction === 'ranked'}
            className="px-4 py-2 rounded-md bg-white text-black hover:bg-zinc-200 font-medium text-xs tracking-wide transition-colors flex items-center gap-2 disabled:opacity-50"
          >
            <Swords className="w-3.5 h-3.5" />
            Quick Ranked Match
          </button>
          <button
            onClick={() => handleStartQueue('casual')}
            disabled={isQueued || loadingAction === 'casual'}
            className="px-4 py-2 rounded-md bg-zinc-900 hover:bg-zinc-800 text-zinc-300 font-medium text-xs border border-zinc-800 transition-colors flex items-center gap-2 disabled:opacity-50"
          >
            <Users className="w-3.5 h-3.5 text-zinc-400" />
            Casual Duel
          </button>
        </div>
      </div>

      {/* Modes Grid */}
      <div className="space-y-4">
        <h2 className="text-sm font-semibold uppercase tracking-wider text-zinc-400">
          Matchmaking Modes
        </h2>

        <div className="grid grid-cols-1 md:grid-cols-3 gap-4">
          {/* Ranked */}
          <div className="p-5 rounded-lg clean-panel clean-card-hover flex flex-col justify-between space-y-4">
            <div className="space-y-2">
              <div className="flex items-center justify-between">
                <span className="text-sm font-semibold text-white">Ranked 1v1</span>
                <span className="text-[10px] font-mono px-1.5 py-0.5 rounded bg-zinc-800 text-zinc-400 border border-zinc-700">
                  ±50 MMR
                </span>
              </div>
              <p className="text-xs text-zinc-400 leading-relaxed">
                Competitive matchmaking. Matches you with peers within your skill window. Rating updates after each battle.
              </p>
            </div>
            <button
              onClick={() => handleStartQueue('ranked')}
              disabled={isQueued}
              className="w-full py-2 px-3 rounded bg-zinc-800 hover:bg-zinc-700 text-zinc-200 text-xs font-medium border border-zinc-700 transition-colors flex items-center justify-center gap-1.5"
            >
              Queue Ranked <ArrowRight className="w-3 h-3 text-zinc-400" />
            </button>
          </div>

          {/* Casual */}
          <div className="p-5 rounded-lg clean-panel clean-card-hover flex flex-col justify-between space-y-4">
            <div className="space-y-2">
              <div className="flex items-center justify-between">
                <span className="text-sm font-semibold text-white">Casual Match</span>
                <span className="text-[10px] font-mono px-1.5 py-0.5 rounded bg-zinc-800 text-zinc-400 border border-zinc-700">
                  Unranked
                </span>
              </div>
              <p className="text-xs text-zinc-400 leading-relaxed">
                Warm up your problem solving without rating pressure. Immediate pairing across any skill bracket.
              </p>
            </div>
            <button
              onClick={() => handleStartQueue('casual')}
              disabled={isQueued}
              className="w-full py-2 px-3 rounded bg-zinc-800 hover:bg-zinc-700 text-zinc-200 text-xs font-medium border border-zinc-700 transition-colors flex items-center justify-center gap-1.5"
            >
              Queue Casual <ArrowRight className="w-3 h-3 text-zinc-400" />
            </button>
          </div>

          {/* Tournament */}
          <div className="p-5 rounded-lg clean-panel clean-card-hover flex flex-col justify-between space-y-4">
            <div className="space-y-2">
              <div className="flex items-center justify-between">
                <span className="text-sm font-semibold text-white">Speed Duel</span>
                <span className="text-[10px] font-mono px-1.5 py-0.5 rounded bg-zinc-800 text-zinc-400 border border-zinc-700">
                  Fast KO
                </span>
              </div>
              <p className="text-xs text-zinc-400 leading-relaxed">
                Short-timer elimination. First engineer to pass all test cases takes the round instantly.
              </p>
            </div>
            <button
              onClick={() => handleStartQueue('tournament')}
              disabled={isQueued}
              className="w-full py-2 px-3 rounded bg-zinc-800 hover:bg-zinc-700 text-zinc-200 text-xs font-medium border border-zinc-700 transition-colors flex items-center justify-center gap-1.5"
            >
              Queue Speed <ArrowRight className="w-3 h-3 text-zinc-400" />
            </button>
          </div>
        </div>
      </div>

      {/* Two Column Layout: Problems & Leaderboard */}
      <div className="grid grid-cols-1 lg:grid-cols-3 gap-8 pt-4">
        {/* Problems List (2 cols) */}
        <div className="lg:col-span-2 space-y-3">
          <div className="flex items-center justify-between pb-2 border-b border-[#27272a]">
            <h3 className="text-sm font-semibold text-white">Challenges</h3>
            <a href="/problems" className="text-xs text-zinc-400 hover:text-white transition-colors flex items-center gap-1">
              All problems <ChevronRight className="w-3 h-3" />
            </a>
          </div>

          <div className="divide-y divide-zinc-800/80 rounded-lg clean-panel overflow-hidden">
            {problems.map((prob) => (
              <div
                key={prob.id}
                className="p-3.5 flex items-center justify-between hover:bg-zinc-900/60 transition-colors"
              >
                <div className="space-y-1">
                  <h4 className="text-xs font-medium text-white">{prob.title}</h4>
                  <div className="text-[11px] text-zinc-500 font-mono">
                    {prob.time_limit_ms}ms limit • {prob.memory_limit_kb / 1024}MB
                  </div>
                </div>

                <span
                  className={`text-[10px] font-mono font-medium uppercase px-2 py-0.5 rounded border ${getDifficultyClass(
                    prob.difficulty
                  )}`}
                >
                  {prob.difficulty}
                </span>
              </div>
            ))}
          </div>
        </div>

        {/* Top Standings (1 col) */}
        <div className="space-y-3">
          <div className="flex items-center justify-between pb-2 border-b border-[#27272a]">
            <h3 className="text-sm font-semibold text-white">Leaderboard</h3>
            <a href="/leaderboard" className="text-xs text-zinc-400 hover:text-white transition-colors flex items-center gap-1">
              Full table <ChevronRight className="w-3 h-3" />
            </a>
          </div>

          <div className="divide-y divide-zinc-800/80 rounded-lg clean-panel overflow-hidden">
            {leaderboard.slice(0, 5).map((player, idx) => (
              <div key={player.id} className="p-3 flex items-center justify-between text-xs">
                <div className="flex items-center gap-3">
                  <span className="font-mono text-zinc-500 w-4 text-center">{idx + 1}</span>
                  <div>
                    <span className="font-medium text-zinc-200 block">{player.name}</span>
                    <span className="text-[11px] text-zinc-500 font-mono">
                      {player.wins}W / {player.losses}L
                    </span>
                  </div>
                </div>
                <span className="font-mono text-zinc-300 font-semibold">{player.rating}</span>
              </div>
            ))}
          </div>
        </div>
      </div>

      <MatchmakingModal />
      <AuthModal isOpen={authModalOpen} onClose={() => setAuthModalOpen(false)} />
    </div>
  );
}
