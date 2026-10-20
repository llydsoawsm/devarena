import type { Metadata } from 'next';
import { Inter, JetBrains_Mono } from 'next/font/google';
import './globals.css';
import { Navbar } from '../components/Navbar';
import { Toaster } from 'sonner';

const inter = Inter({
  variable: '--font-sans',
  subsets: ['latin'],
});

const jetbrainsMono = JetBrains_Mono({
  variable: '--font-mono',
  subsets: ['latin'],
});

export const metadata: Metadata = {
  title: 'DevCode - Real-Time 1v1 Competitive Multiplayer Coding Arena',
  description:
    'Head-to-head real-time algorithmic coding arena with isolated sandboxed execution, Monaco Editor, live matchmaking, and global ELO ladder.',
};

export default function RootLayout({
  children,
}: Readonly<{
  children: React.ReactNode;
}>) {
  return (
    <html lang="en" className={`${inter.variable} ${jetbrainsMono.variable} dark`}>
      <body className="min-h-screen bg-[#07090e] text-slate-100 flex flex-col font-sans selection:bg-indigo-500/30 selection:text-cyan-300">
        <Navbar />
        <main className="flex-1">{children}</main>
        <Toaster position="bottom-right" theme="dark" richColors />
      </body>
    </html>
  );
}
