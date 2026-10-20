'use client';

import React, { useEffect, useState } from 'react';
import { matchesApi } from '../../lib/api';
import { User } from '../../lib/types';

export default function LeaderboardPage() {
  const [leaderboard, setLeaderboard] = useState<User[]>([]);
  const [loading, setLoading] = useState(true);

  useEffect(() => {
    matchesApi.leaderboard().then((res) => {
      setLeaderboard(res.data.leaderboard || []);
      setLoading(false);
    });
  }, []);

  return (
    <div className="min-h-screen py-10 px-4 sm:px-6 max-w-6xl mx-auto space-y-6">
      <div className="space-y-1">
        <h1 className="text-2xl font-bold text-white tracking-tight">Arena Leaderboard</h1>
        <p className="text-xs text-zinc-400">Global rankings calculated via standard logistic ELO.</p>
      </div>

      <div className="rounded-lg clean-panel overflow-hidden border border-[#27272a]">
        <div className="overflow-x-auto">
          <table className="w-full text-left text-xs">
            <thead className="bg-[#16161a] text-zinc-400 uppercase tracking-wider font-medium text-[11px] border-b border-[#27272a]">
              <tr>
                <th className="py-3 px-4 w-16">Rank</th>
                <th className="py-3 px-4">Duelist</th>
                <th className="py-3 px-4">Rating</th>
                <th className="py-3 px-4">Win Rate</th>
                <th className="py-3 px-4 text-right">Matches</th>
              </tr>
            </thead>
            <tbody className="divide-y divide-zinc-800/80">
              {loading ? (
                <tr>
                  <td colSpan={5} className="py-8 text-center text-zinc-500">
                    Loading leaderboard...
                  </td>
                </tr>
              ) : leaderboard.length > 0 ? (
                leaderboard.map((user, idx) => {
                  const winRate =
                    user.matches_played > 0
                      ? Math.round((user.wins / user.matches_played) * 100)
                      : 0;

                  return (
                    <tr key={user.id} className="hover:bg-zinc-900/50 transition-colors">
                      <td className="py-3.5 px-4 font-mono text-zinc-500">{idx + 1}</td>
                      <td className="py-3.5 px-4 font-medium text-white">{user.name}</td>
                      <td className="py-3.5 px-4 font-mono text-zinc-200 font-semibold">
                        {user.rating} ELO
                      </td>
                      <td className="py-3.5 px-4 font-mono text-zinc-400">
                        {winRate}% ({user.wins}W / {user.losses}L)
                      </td>
                      <td className="py-3.5 px-4 text-right font-mono text-zinc-500">
                        {user.matches_played}
                      </td>
                    </tr>
                  );
                })
              ) : (
                <tr>
                  <td colSpan={5} className="py-8 text-center text-zinc-500">
                    No duels recorded yet.
                  </td>
                </tr>
              )}
            </tbody>
          </table>
        </div>
      </div>
    </div>
  );
}
