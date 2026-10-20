'use client';

import React, { useState } from 'react';
import { TestCase, ExecutionResult } from '../lib/types';

interface ConsoleDrawerProps {
  testCases: TestCase[];
  executionResult: ExecutionResult | null;
  isRunning: boolean;
}

export function ConsoleDrawer({ testCases, executionResult, isRunning }: ConsoleDrawerProps) {
  const [activeTab, setActiveTab] = useState<'tests' | 'result'>('tests');
  const [selectedTestCaseIdx, setSelectedTestCaseIdx] = useState(0);

  const selectedCase = testCases[selectedTestCaseIdx];

  const getVerdictBadge = (verdict: string) => {
    switch (verdict) {
      case 'AC':
        return (
          <span className="px-2 py-0.5 rounded text-[11px] font-mono font-medium text-emerald-400 bg-emerald-950/40 border border-emerald-800/50">
            Accepted
          </span>
        );
      case 'WA':
        return (
          <span className="px-2 py-0.5 rounded text-[11px] font-mono font-medium text-rose-400 bg-rose-950/40 border border-rose-800/50">
            Wrong Answer
          </span>
        );
      case 'TLE':
        return (
          <span className="px-2 py-0.5 rounded text-[11px] font-mono font-medium text-amber-400 bg-amber-950/40 border border-amber-800/50">
            Time Limit Exceeded
          </span>
        );
      default:
        return (
          <span className="px-2 py-0.5 rounded text-[11px] font-mono font-medium text-zinc-300 bg-zinc-800 border border-zinc-700">
            Error ({verdict})
          </span>
        );
    }
  };

  return (
    <div className="flex flex-col h-full bg-[#0d0d10]">
      {/* Console Nav Tabs */}
      <div className="flex items-center justify-between px-3 py-1.5 bg-[#121215] border-b border-[#27272a] text-xs">
        <div className="flex items-center gap-1">
          <button
            onClick={() => setActiveTab('tests')}
            className={`px-2.5 py-1 rounded text-xs font-medium transition-colors ${
              activeTab === 'tests'
                ? 'bg-zinc-800 text-white'
                : 'text-zinc-400 hover:text-white'
            }`}
          >
            Test Cases
          </button>
          <button
            onClick={() => setActiveTab('result')}
            className={`px-2.5 py-1 rounded text-xs font-medium transition-colors ${
              activeTab === 'result'
                ? 'bg-zinc-800 text-white'
                : 'text-zinc-400 hover:text-white'
            }`}
          >
            Execution Output
          </button>
        </div>

        {executionResult && (
          <div className="flex items-center gap-2">
            {getVerdictBadge(executionResult.verdict)}
            <span className="text-[11px] font-mono text-zinc-500">
              {executionResult.passed_test_cases}/{executionResult.total_test_cases} passed • {executionResult.runtime_ms}ms
            </span>
          </div>
        )}
      </div>

      {/* Console Content */}
      <div className="flex-1 p-3.5 overflow-y-auto font-mono text-xs text-zinc-300">
        {isRunning ? (
          <div className="flex items-center justify-center h-full gap-2 text-zinc-500">
            <span className="w-3.5 h-3.5 border-2 border-zinc-600 border-t-white rounded-full animate-spin" />
            <span>Executing test cases...</span>
          </div>
        ) : activeTab === 'tests' ? (
          <div className="space-y-3">
            <div className="flex items-center gap-1.5">
              {testCases.map((tc, idx) => (
                <button
                  key={tc.id || idx}
                  onClick={() => setSelectedTestCaseIdx(idx)}
                  className={`px-2 py-0.5 rounded text-[11px] font-mono transition-colors ${
                    selectedTestCaseIdx === idx
                      ? 'bg-zinc-700 text-white'
                      : 'bg-zinc-900 text-zinc-500 hover:text-zinc-300'
                  }`}
                >
                  Case {idx + 1}
                </button>
              ))}
            </div>

            {selectedCase ? (
              <div className="grid grid-cols-2 gap-3 text-[11px]">
                <div className="space-y-1">
                  <span className="text-zinc-500 uppercase text-[10px] block">Input</span>
                  <pre className="p-2.5 rounded bg-zinc-950 border border-zinc-800/80 text-zinc-300 whitespace-pre-wrap">
                    {selectedCase.input || '<empty>'}
                  </pre>
                </div>
                <div className="space-y-1">
                  <span className="text-zinc-500 uppercase text-[10px] block">Expected Output</span>
                  <pre className="p-2.5 rounded bg-zinc-950 border border-zinc-800/80 text-zinc-300 whitespace-pre-wrap">
                    {selectedCase.expected_output || '<empty>'}
                  </pre>
                </div>
              </div>
            ) : (
              <p className="text-zinc-600">No test cases available.</p>
            )}
          </div>
        ) : (
          <div>
            {executionResult ? (
              <div className="space-y-3">
                <div className="flex items-center gap-5 p-2.5 rounded bg-zinc-950 border border-zinc-800/80 text-xs">
                  <div>
                    <span className="text-[10px] text-zinc-500 block uppercase">Result</span>
                    <span className="font-semibold text-white">{executionResult.verdict}</span>
                  </div>
                  <div>
                    <span className="text-[10px] text-zinc-500 block uppercase">Passed</span>
                    <span className="text-zinc-300">
                      {executionResult.passed_test_cases}/{executionResult.total_test_cases}
                    </span>
                  </div>
                  <div>
                    <span className="text-[10px] text-zinc-500 block uppercase">Runtime</span>
                    <span className="text-zinc-300">{executionResult.runtime_ms}ms</span>
                  </div>
                </div>

                {executionResult.error_details && (
                  <div className="space-y-1">
                    <span className="text-[10px] uppercase text-zinc-500 block">Diff / Error</span>
                    <pre className="p-2.5 rounded bg-zinc-950 border border-zinc-800/80 text-rose-300 whitespace-pre-wrap">
                      {executionResult.error_details}
                    </pre>
                  </div>
                )}
              </div>
            ) : (
              <p className="text-zinc-600">Press Run (Ctrl+Enter) to execute code.</p>
            )}
          </div>
        )}
      </div>
    </div>
  );
}
