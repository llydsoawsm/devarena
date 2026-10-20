'use client';

import React, { useEffect, useState, Suspense } from 'react';
import { useRouter, useParams } from 'next/navigation';
import { useAuthStore } from '../../../stores/authStore';
import { useMatchStore } from '../../../stores/matchStore';
import { matchesApi, submissionsApi } from '../../../lib/api';
import { Match, ExecutionResult } from '../../../lib/types';
import { useDevCodeSocket } from '../../../hooks/useDevCodeSocket';
import { MonacoArena } from '../../../components/MonacoArena';
import { ConsoleDrawer } from '../../../components/ConsoleDrawer';
import { VictoryModal } from '../../../components/VictoryModal';
import { Clock, Flag } from 'lucide-react';
import { toast } from 'sonner';

function MatchArenaContent() {
  const router = useRouter();
  const routeParams = useParams();
  const matchPublicId = (routeParams?.id as string) || '';

  const { user, loadUser } = useAuthStore();
  const {
    activeMatch,
    setActiveMatch,
    countdown,
    isMatchLive,
    opponentProgress,
  } = useMatchStore();

  const [match, setMatch] = useState<Match | null>(null);
  const [loading, setLoading] = useState(true);
  const [isRunningSample, setIsRunningSample] = useState(false);
  const [isSubmitting, setIsSubmitting] = useState(false);
  const [executionResult, setExecutionResult] = useState<ExecutionResult | null>(null);
  const [matchTimer, setMatchTimer] = useState<number>(1800);

  const { sendParticipantState } = useDevCodeSocket(matchPublicId);

  useEffect(() => {
    loadUser();

    matchesApi
      .getByPublicId(matchPublicId)
      .then((res) => {
        setMatch(res.data.match);
        setActiveMatch(res.data.match);
        setLoading(false);
      })
      .catch(() => {
        toast.error('Match not found');
        router.push('/');
      });
  }, [matchPublicId, loadUser, setActiveMatch, router]);

  useEffect(() => {
    let timerInterval: NodeJS.Timeout;
    if (isMatchLive && matchTimer > 0) {
      timerInterval = setInterval(() => {
        setMatchTimer((prev) => (prev > 0 ? prev - 1 : 0));
      }, 1000);
    }
    return () => clearInterval(timerInterval);
  }, [isMatchLive, matchTimer]);

  const formatTimer = (seconds: number) => {
    const mins = Math.floor(seconds / 60);
    const secs = seconds % 60;
    return `${mins.toString().padStart(2, '0')}:${secs.toString().padStart(2, '0')}`;
  };

  const opponentParticipant = match?.participants?.find((p) => p.user_id !== user?.id);

  const handleRunCode = async (language: string, code: string) => {
    if (!match?.problem_id) return;
    setIsRunningSample(true);
    try {
      const res = await submissionsApi.runSample({
        problem_id: match.problem_id,
        language,
        code,
      });
      setExecutionResult(res.data.result);
      if (res.data.result.verdict === 'AC') {
        toast.success('Sample tests passed');
      } else {
        toast.error(`Sample tests failed: ${res.data.result.verdict}`);
      }
    } catch {
      toast.error('Failed to run code in sandbox');
    } finally {
      setIsRunningSample(false);
    }
  };

  const handleSubmitCode = async (language: string, code: string) => {
    setIsSubmitting(true);
    try {
      const res = await submissionsApi.submitMatch({
        match_public_id: matchPublicId,
        language,
        code,
      });

      setExecutionResult(res.data.result);
      sendParticipantState({
        passed_tests: res.data.result.passed_test_cases,
        total_tests: res.data.result.total_test_cases,
        verdict: res.data.result.verdict,
      });

      if (res.data.result.verdict === 'AC') {
        toast.success('Accepted! All tests passed.');
      } else {
        toast.warning(
          `Tested: ${res.data.result.passed_test_cases}/${res.data.result.total_test_cases} passed (${res.data.result.verdict})`
        );
      }
    } catch (err: any) {
      toast.error(err.response?.data?.error || 'Failed to submit solution');
    } finally {
      setIsSubmitting(false);
    }
  };

  const handleForfeit = async () => {
    if (confirm('Forfeit this match? Your rating will decrease.')) {
      try {
        await matchesApi.forfeit(matchPublicId);
        toast.info('Match forfeited');
      } catch {
        toast.error('Forfeit failed');
      }
    }
  };

  if (loading || !match) {
    return (
      <div className="flex flex-col items-center justify-center min-h-[70vh] gap-2 text-zinc-500 text-xs">
        <span className="w-4 h-4 border-2 border-zinc-600 border-t-white rounded-full animate-spin" />
        <span>Loading match arena...</span>
      </div>
    );
  }

  const starterCodes =
    typeof match.problem?.starter_code === 'string'
      ? JSON.parse(match.problem.starter_code)
      : match.problem?.starter_code || {};

  return (
    <div className="flex flex-col h-[calc(100vh-3.5rem)] overflow-hidden bg-[#09090b]">
      {/* Top Arena HUD Bar */}
      <div className="px-5 py-2.5 bg-[#0e0e11] border-b border-[#27272a] flex items-center justify-between text-xs">
        {/* Left: Problem & Difficulty */}
        <div className="flex items-center gap-3">
          <span className="font-semibold text-white tracking-tight">{match.problem?.title}</span>
          <span className="text-[10px] uppercase font-mono px-1.5 py-0.5 rounded bg-zinc-800 text-zinc-400 border border-zinc-700">
            {match.mode}
          </span>
          <span className="text-[10px] uppercase font-mono px-1.5 py-0.5 rounded bg-zinc-800 text-zinc-400 border border-zinc-700">
            {match.problem?.difficulty}
          </span>
        </div>

        {/* Center: Live Timer Clock */}
        <div className="flex items-center gap-1.5 font-mono text-xs font-medium text-zinc-300">
          <Clock className="w-3.5 h-3.5 text-zinc-400" />
          <span>{formatTimer(matchTimer)}</span>
        </div>

        {/* Right: Opponent Progress & Forfeit */}
        <div className="flex items-center gap-5">
          {opponentParticipant && (
            <div className="flex items-center gap-2 text-zinc-400">
              <span>{opponentParticipant.user?.name}</span>
              <span className="font-mono text-[11px] px-1.5 py-0.5 rounded bg-zinc-900 border border-zinc-800 text-zinc-300">
                {opponentProgress ? `${opponentProgress.passed_tests}/${opponentProgress.total_tests}` : '0/4'} tests
              </span>
            </div>
          )}

          <button
            onClick={handleForfeit}
            className="p-1 rounded text-zinc-500 hover:text-rose-400 transition-colors"
            title="Forfeit"
          >
            <Flag className="w-3.5 h-3.5" />
          </button>
        </div>
      </div>

      {/* Main Split Screen Arena */}
      <div className="flex-1 grid grid-cols-1 lg:grid-cols-12 overflow-hidden">
        {/* Left: Problem Statement (5 cols) */}
        <div className="lg:col-span-5 h-full overflow-y-auto p-6 border-r border-[#27272a] bg-[#09090b] text-zinc-300 space-y-6">
          <div className="space-y-3">
            <h3 className="text-lg font-semibold text-white">{match.problem?.title}</h3>
            <div className="text-xs leading-relaxed text-zinc-400 whitespace-pre-wrap font-sans">
              {match.problem?.description}
            </div>
          </div>

          <div className="p-3.5 rounded border border-[#27272a] bg-zinc-900/40 text-[11px] font-mono text-zinc-400 space-y-1">
            <div className="text-zinc-500 uppercase font-semibold text-[10px]">Limits</div>
            <div>Time Limit: {match.problem?.time_limit_ms}ms</div>
            <div>Memory Limit: {match.problem?.memory_limit_kb / 1024}MB</div>
          </div>
        </div>

        {/* Right: Editor & Console (7 cols) */}
        <div className="lg:col-span-7 flex flex-col h-full overflow-hidden bg-[#0c0c0e]">
          <div className="flex-[65] overflow-hidden p-2">
            <MonacoArena
              starterCodes={starterCodes}
              onRunCode={handleRunCode}
              onSubmitCode={handleSubmitCode}
              isLocked={countdown > 0}
              countdown={countdown}
              isRunning={isRunningSample}
              isSubmitting={isSubmitting}
            />
          </div>

          <div className="flex-[35] overflow-hidden border-t border-[#27272a]">
            <ConsoleDrawer
              testCases={match.problem?.test_cases || []}
              executionResult={executionResult}
              isRunning={isRunningSample || isSubmitting}
            />
          </div>
        </div>
      </div>

      <VictoryModal />
    </div>
  );
}

export default function MatchArenaPage() {
  return (
    <Suspense
      fallback={
        <div className="flex flex-col items-center justify-center min-h-[70vh] gap-2 text-zinc-500 text-xs">
          <span className="w-4 h-4 border-2 border-zinc-600 border-t-white rounded-full animate-spin" />
          <span>Loading match arena...</span>
        </div>
      }
    >
      <MatchArenaContent />
    </Suspense>
  );
}
