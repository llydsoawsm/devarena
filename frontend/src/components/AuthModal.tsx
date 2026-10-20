'use client';

import React, { useState } from 'react';
import { useAuthStore } from '../stores/authStore';
import { X } from 'lucide-react';
import { toast } from 'sonner';

interface AuthModalProps {
  isOpen: boolean;
  onClose: () => void;
}

export function AuthModal({ isOpen, onClose }: AuthModalProps) {
  const [isRegister, setIsRegister] = useState(false);
  const [name, setName] = useState('');
  const [email, setEmail] = useState('');
  const [password, setPassword] = useState('');
  const [loading, setLoading] = useState(false);

  const { login, register } = useAuthStore();

  if (!isOpen) return null;

  const handleSubmit = async (e: React.FormEvent) => {
    e.preventDefault();
    setLoading(true);
    try {
      if (isRegister) {
        await register(name, email, password);
        toast.success(`Welcome, ${name}`);
      } else {
        await login(email, password);
        toast.success('Logged in');
      }
      onClose();
    } catch (err: any) {
      toast.error(err.response?.data?.error || 'Authentication failed');
    } finally {
      setLoading(false);
    }
  };

  const handleQuickPlay = async (playerNum: 1 | 2) => {
    setLoading(true);
    const demoEmail = `player${playerNum}@devcode.io`;
    const demoPass = 'password123';
    const demoName = `Warrior_${playerNum}`;

    try {
      await login(demoEmail, demoPass);
      toast.success(`Logged in as ${demoName}`);
      onClose();
    } catch {
      try {
        await register(demoName, demoEmail, demoPass);
        toast.success(`Logged in as ${demoName}`);
        onClose();
      } catch {
        toast.error('Quick play failed');
      }
    } finally {
      setLoading(false);
    }
  };

  return (
    <div className="fixed inset-0 z-50 flex items-center justify-center p-4 bg-black/60 backdrop-blur-xs">
      <div className="relative w-full max-w-sm p-6 rounded-lg bg-[#121215] border border-[#27272a] shadow-xl text-left space-y-5">
        <button
          onClick={onClose}
          className="absolute top-4 right-4 text-zinc-500 hover:text-zinc-200 transition-colors"
        >
          <X className="w-4 h-4" />
        </button>

        <div>
          <h2 className="text-base font-semibold text-white">
            {isRegister ? 'Create Account' : 'Sign In'}
          </h2>
          <p className="text-xs text-zinc-400 mt-0.5">
            Join the real-time coding arena
          </p>
        </div>

        {/* 1-Click Test Quick Login */}
        <div className="p-3 rounded bg-zinc-900 border border-zinc-800 space-y-2">
          <span className="text-[10px] uppercase font-semibold text-zinc-400 tracking-wider block">
            1-Click Test Accounts
          </span>
          <div className="grid grid-cols-2 gap-2">
            <button
              type="button"
              onClick={() => handleQuickPlay(1)}
              disabled={loading}
              className="py-1 px-2.5 text-xs font-medium rounded bg-zinc-800 hover:bg-zinc-700 text-zinc-300 border border-zinc-700 transition-colors"
            >
              Player 1
            </button>
            <button
              type="button"
              onClick={() => handleQuickPlay(2)}
              disabled={loading}
              className="py-1 px-2.5 text-xs font-medium rounded bg-zinc-800 hover:bg-zinc-700 text-zinc-300 border border-zinc-700 transition-colors"
            >
              Player 2
            </button>
          </div>
        </div>

        <form onSubmit={handleSubmit} className="space-y-3 text-xs">
          {isRegister && (
            <div className="space-y-1">
              <label className="text-zinc-400">Name</label>
              <input
                type="text"
                required
                value={name}
                onChange={(e) => setName(e.target.value)}
                placeholder="Your name"
                className="w-full px-3 py-1.5 rounded bg-zinc-900 border border-zinc-800 text-white placeholder-zinc-600 focus:outline-none focus:border-zinc-500"
              />
            </div>
          )}

          <div className="space-y-1">
            <label className="text-zinc-400">Email</label>
            <input
              type="email"
              required
              value={email}
              onChange={(e) => setEmail(e.target.value)}
              placeholder="name@domain.com"
              className="w-full px-3 py-1.5 rounded bg-zinc-900 border border-zinc-800 text-white placeholder-zinc-600 focus:outline-none focus:border-zinc-500"
            />
          </div>

          <div className="space-y-1">
            <label className="text-zinc-400">Password</label>
            <input
              type="password"
              required
              value={password}
              onChange={(e) => setPassword(e.target.value)}
              placeholder="••••••••"
              className="w-full px-3 py-1.5 rounded bg-zinc-900 border border-zinc-800 text-white placeholder-zinc-600 focus:outline-none focus:border-zinc-500"
            />
          </div>

          <button
            type="submit"
            disabled={loading}
            className="w-full mt-2 py-2 px-3 rounded bg-white hover:bg-zinc-200 text-black font-medium text-xs transition-colors flex items-center justify-center disabled:opacity-50"
          >
            {loading ? 'Please wait...' : isRegister ? 'Register' : 'Sign In'}
          </button>
        </form>

        <div className="text-center pt-1 border-t border-zinc-800">
          <button
            type="button"
            onClick={() => setIsRegister(!isRegister)}
            className="text-[11px] text-zinc-500 hover:text-zinc-300 transition-colors"
          >
            {isRegister ? 'Already have an account? Sign In' : "Don't have an account? Register"}
          </button>
        </div>
      </div>
    </div>
  );
}
