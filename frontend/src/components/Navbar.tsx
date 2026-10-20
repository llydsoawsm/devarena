'use client';

import React, { useState } from 'react';
import { useAuthStore } from '../stores/authStore';
import { AuthModal } from './AuthModal';
import { Terminal, LogOut, User as UserIcon } from 'lucide-react';

export function Navbar() {
  const { user, logout } = useAuthStore();
  const [authModalOpen, setAuthModalOpen] = useState(false);

  return (
    <>
      <header className="sticky top-0 z-40 w-full bg-[#09090b]/90 backdrop-blur-md border-b border-[#27272a]">
        <div className="max-w-6xl mx-auto px-4 sm:px-6 h-14 flex items-center justify-between">
          {/* Logo & Navigation */}
          <div className="flex items-center gap-8">
            <a href="/" className="flex items-center gap-2 text-white font-semibold tracking-tight text-sm">
              <div className="w-6 h-6 rounded bg-zinc-800 border border-zinc-700 flex items-center justify-center text-zinc-300">
                <Terminal className="w-3.5 h-3.5" />
              </div>
              <span>devcode</span>
            </a>

            <nav className="hidden sm:flex items-center gap-5 text-xs font-medium">
              <a href="/" className="text-zinc-300 hover:text-white transition-colors">
                Arena
              </a>
              <a href="/problems" className="text-zinc-400 hover:text-white transition-colors">
                Problems
              </a>
              <a href="/leaderboard" className="text-zinc-400 hover:text-white transition-colors">
                Leaderboard
              </a>
            </nav>
          </div>

          {/* User Profile / Auth */}
          <div className="flex items-center gap-3">
            {user ? (
              <div className="flex items-center gap-3">
                <div className="flex items-center gap-2 text-xs">
                  <span className="font-medium text-zinc-200">{user.name}</span>
                  <span className="px-1.5 py-0.5 rounded bg-zinc-800 text-zinc-400 font-mono text-[11px] border border-zinc-700">
                    {user.rating} ELO
                  </span>
                  <span className="text-[11px] text-zinc-500 font-mono hidden sm:inline">
                    ({user.wins}W - {user.losses}L)
                  </span>
                </div>

                <button
                  onClick={logout}
                  title="Sign Out"
                  className="p-1.5 rounded text-zinc-400 hover:text-zinc-200 hover:bg-zinc-800 transition-colors"
                >
                  <LogOut className="w-3.5 h-3.5" />
                </button>
              </div>
            ) : (
              <button
                onClick={() => setAuthModalOpen(true)}
                className="px-3 py-1.5 rounded-md bg-white text-black hover:bg-zinc-200 font-medium text-xs transition-colors flex items-center gap-1.5"
              >
                <UserIcon className="w-3.5 h-3.5" />
                Sign In
              </button>
            )}
          </div>
        </div>
      </header>

      <AuthModal isOpen={authModalOpen} onClose={() => setAuthModalOpen(false)} />
    </>
  );
}
