'use client';

import React, { useState, useEffect } from 'react';
import Editor from '@monaco-editor/react';
import { Play, Send, RotateCcw } from 'lucide-react';

interface MonacoArenaProps {
  starterCodes: Record<string, string>;
  onRunCode: (language: string, code: string) => void;
  onSubmitCode: (language: string, code: string) => void;
  isLocked: boolean;
  countdown: number;
  isRunning: boolean;
  isSubmitting: boolean;
}

export function MonacoArena({
  starterCodes,
  onRunCode,
  onSubmitCode,
  isLocked,
  countdown,
  isRunning,
  isSubmitting,
}: MonacoArenaProps) {
  const [language, setLanguage] = useState<string>('python');
  const [code, setCode] = useState<string>('');

  useEffect(() => {
    if (starterCodes && starterCodes[language]) {
      setCode(starterCodes[language]);
    }
  }, [language, starterCodes]);

  const handleLanguageChange = (newLang: string) => {
    setLanguage(newLang);
    if (starterCodes && starterCodes[newLang]) {
      setCode(starterCodes[newLang]);
    }
  };

  const handleReset = () => {
    if (starterCodes && starterCodes[language]) {
      setCode(starterCodes[language]);
    }
  };

  useEffect(() => {
    const handleKeyDown = (e: KeyboardEvent) => {
      if ((e.ctrlKey || e.metaKey) && e.key === 'Enter') {
        e.preventDefault();
        if (e.shiftKey) {
          if (!isLocked && !isSubmitting) onSubmitCode(language, code);
        } else {
          if (!isLocked && !isRunning) onRunCode(language, code);
        }
      }
    };
    window.addEventListener('keydown', handleKeyDown);
    return () => window.removeEventListener('keydown', handleKeyDown);
  }, [isLocked, isRunning, isSubmitting, language, code, onRunCode, onSubmitCode]);

  return (
    <div className="relative flex flex-col h-full bg-[#121215] rounded-md overflow-hidden border border-[#27272a]">
      {/* Editor Header Toolbar */}
      <div className="flex items-center justify-between px-3 py-2 bg-[#16161a] border-b border-[#27272a] text-xs">
        <div className="flex items-center gap-2">
          <select
            value={language}
            onChange={(e) => handleLanguageChange(e.target.value)}
            disabled={isLocked}
            className="px-2 py-1 text-xs rounded bg-zinc-800 text-zinc-200 border border-zinc-700 focus:outline-none cursor-pointer disabled:opacity-50"
          >
            <option value="python">Python</option>
            <option value="javascript">JavaScript</option>
            <option value="typescript">TypeScript</option>
            <option value="go">Go</option>
          </select>

          <button
            onClick={handleReset}
            disabled={isLocked}
            title="Reset to starter template"
            className="p-1 rounded text-zinc-400 hover:text-white transition-colors disabled:opacity-50"
          >
            <RotateCcw className="w-3.5 h-3.5" />
          </button>
        </div>

        <div className="flex items-center gap-2">
          <button
            onClick={() => onRunCode(language, code)}
            disabled={isLocked || isRunning}
            className="px-2.5 py-1 rounded bg-zinc-800 hover:bg-zinc-700 text-zinc-200 font-medium text-xs border border-zinc-700 transition-colors flex items-center gap-1.5 disabled:opacity-50"
          >
            {isRunning ? (
              <span className="w-3 h-3 border-2 border-zinc-400 border-t-white rounded-full animate-spin" />
            ) : (
              <Play className="w-3 h-3 text-zinc-400" />
            )}
            Run <span className="text-[10px] text-zinc-500 font-mono hidden sm:inline">(Ctrl+Enter)</span>
          </button>

          <button
            onClick={() => onSubmitCode(language, code)}
            disabled={isLocked || isSubmitting}
            className="px-3 py-1 rounded bg-white hover:bg-zinc-200 text-black font-medium text-xs transition-colors flex items-center gap-1.5 disabled:opacity-50"
          >
            {isSubmitting ? (
              <span className="w-3 h-3 border-2 border-zinc-600 border-t-black rounded-full animate-spin" />
            ) : (
              <Send className="w-3 h-3" />
            )}
            Submit
          </button>
        </div>
      </div>

      {/* Editor Content */}
      <div className="relative flex-1">
        <Editor
          height="100%"
          language={language === 'typescript' ? 'typescript' : language === 'javascript' ? 'javascript' : language === 'go' ? 'go' : 'python'}
          value={code}
          onChange={(val) => setCode(val || '')}
          theme="vs-dark"
          options={{
            readOnly: isLocked,
            minimap: { enabled: false },
            fontSize: 13,
            lineNumbers: 'on',
            scrollBeyondLastLine: false,
            automaticLayout: true,
            tabSize: 4,
            cursorBlinking: 'smooth',
            padding: { top: 10, bottom: 10 },
          }}
        />

        {/* Minimalist Countdown Screen */}
        {isLocked && countdown > 0 && (
          <div className="absolute inset-0 z-20 flex flex-col items-center justify-center bg-[#09090b]/90 backdrop-blur-sm text-center p-6">
            <span className="text-xs uppercase tracking-widest text-zinc-500 font-semibold mb-2">
              Unlocks in
            </span>
            <div className="text-6xl font-mono font-bold text-white tracking-tight">
              {countdown}
            </div>
          </div>
        )}
      </div>
    </div>
  );
}
