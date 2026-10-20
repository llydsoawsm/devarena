'use client';

import React, { useEffect, useState } from 'react';
import { problemsApi } from '../../lib/api';
import { Problem } from '../../lib/types';
import { Search } from 'lucide-react';

export default function ProblemsPage() {
  const [problems, setProblems] = useState<Problem[]>([]);
  const [search, setSearch] = useState('');
  const [difficultyFilter, setDifficultyFilter] = useState<'all' | 'easy' | 'medium' | 'hard'>('all');
  const [loading, setLoading] = useState(true);

  useEffect(() => {
    problemsApi.list().then((res) => {
      setProblems(res.data.problems || []);
      setLoading(false);
    });
  }, []);

  const filteredProblems = problems.filter((p) => {
    const matchesSearch = p.title.toLowerCase().includes(search.toLowerCase());
    const matchesDiff = difficultyFilter === 'all' || p.difficulty === difficultyFilter;
    return matchesSearch && matchesDiff;
  });

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
    <div className="min-h-screen py-10 px-4 sm:px-6 max-w-6xl mx-auto space-y-6">
      <div className="space-y-1">
        <h1 className="text-2xl font-bold text-white tracking-tight">Problem Vault</h1>
        <p className="text-xs text-zinc-400">Algorithmic problems available for 1v1 arenas and practice.</p>
      </div>

      {/* Filter and Search Bar */}
      <div className="flex flex-col sm:flex-row gap-3 items-center justify-between pb-4 border-b border-[#27272a]">
        <div className="relative w-full sm:w-80">
          <Search className="absolute left-3 top-1/2 -translate-y-1/2 w-3.5 h-3.5 text-zinc-500" />
          <input
            type="text"
            value={search}
            onChange={(e) => setSearch(e.target.value)}
            placeholder="Search problems..."
            className="w-full pl-9 pr-3 py-1.5 text-xs rounded bg-zinc-900 border border-zinc-800 text-white placeholder-zinc-500 focus:outline-none focus:border-zinc-600"
          />
        </div>

        <div className="flex items-center gap-1.5 w-full sm:w-auto">
          {(['all', 'easy', 'medium', 'hard'] as const).map((diff) => (
            <button
              key={diff}
              onClick={() => setDifficultyFilter(diff)}
              className={`px-2.5 py-1 rounded text-xs font-medium capitalize transition-colors ${
                difficultyFilter === diff
                  ? 'bg-zinc-800 text-white border border-zinc-700'
                  : 'text-zinc-400 hover:text-white'
              }`}
            >
              {diff}
            </button>
          ))}
        </div>
      </div>

      {/* Problems List */}
      <div className="divide-y divide-zinc-800/80 rounded-lg clean-panel overflow-hidden">
        {loading ? (
          <div className="p-8 text-center text-xs text-zinc-500">Loading problems...</div>
        ) : filteredProblems.length > 0 ? (
          filteredProblems.map((prob) => (
            <div
              key={prob.id}
              className="p-4 flex items-center justify-between hover:bg-zinc-900/50 transition-colors"
            >
              <div className="space-y-1">
                <h3 className="text-sm font-medium text-white">{prob.title}</h3>
                <div className="text-[11px] text-zinc-500 font-mono">
                  Time: {prob.time_limit_ms}ms • RAM: {prob.memory_limit_kb / 1024}MB
                </div>
              </div>

              <div className="flex items-center gap-3">
                <span
                  className={`text-[10px] font-mono font-medium uppercase px-2 py-0.5 rounded border ${getDifficultyClass(
                    prob.difficulty
                  )}`}
                >
                  {prob.difficulty}
                </span>
                <a
                  href="/"
                  className="px-2.5 py-1 rounded bg-zinc-800 hover:bg-zinc-700 text-zinc-300 text-xs font-medium border border-zinc-700 transition-colors"
                >
                  Play
                </a>
              </div>
            </div>
          ))
        ) : (
          <div className="p-8 text-center text-xs text-zinc-500">No problems found.</div>
        )}
      </div>
    </div>
  );
}
