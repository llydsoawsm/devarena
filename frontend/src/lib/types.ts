export interface User {
  id: number;
  name: string;
  email: string;
  rating: number;
  matches_played: number;
  wins: number;
  losses: number;
  created_at: string;
}

export interface TestCase {
  id: number;
  problem_id: number;
  input: string;
  expected_output: string;
  is_hidden: boolean;
  sequence_order: number;
}

export interface Problem {
  id: number;
  slug: string;
  title: string;
  difficulty: 'easy' | 'medium' | 'hard';
  description: string;
  starter_code: Record<string, string>;
  time_limit_ms: number;
  memory_limit_kb: number;
  test_cases?: TestCase[];
}

export interface MatchParticipant {
  id: number;
  match_id: number;
  user_id: number;
  user: User;
  seat: 1 | 2;
  result: 'pending' | 'won' | 'lost' | 'draw' | 'forfeit';
  rating_before: number;
  rating_after: number;
  rating_delta: number;
  submission_count: number;
  passed_tests: number;
  total_tests: number;
  is_connected: boolean;
}

export interface Match {
  id: number;
  public_id: string;
  mode: 'ranked' | 'casual' | 'tournament';
  status: 'created' | 'countdown' | 'live' | 'finished' | 'settled' | 'cancelled';
  problem_id: number;
  problem: Problem;
  max_duration_seconds: number;
  countdown_seconds: number;
  started_at?: string;
  ended_at?: string;
  settled_at?: string;
  winner_id?: number;
  participants: MatchParticipant[];
}

export interface ExecutionResult {
  verdict: 'AC' | 'WA' | 'TLE' | 'MLE' | 'CE' | 'RE';
  passed_test_cases: number;
  total_test_cases: number;
  runtime_ms: number;
  memory_kb: number;
  compile_output?: string;
  error_details?: string;
}

export interface Submission {
  id: number;
  match_id?: number;
  user_id: number;
  problem_id: number;
  language: string;
  code: string;
  verdict: 'AC' | 'WA' | 'TLE' | 'MLE' | 'CE' | 'RE';
  passed_test_cases: number;
  total_test_cases: number;
  runtime_ms: number;
  memory_kb: number;
  compile_output?: string;
  submitted_at: string;
}

export interface OpponentState {
  user_id: number;
  passed_tests: number;
  total_tests: number;
  verdict?: string;
}
