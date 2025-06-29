"use client";
import useAuthStore from '@/stores/useAuthStore';
import React, { useState } from 'react';

const JoinGameOverlay: React.FC = () => {
  const [username, setUsername] = useState('');
    const {user, joinGame} = useAuthStore()

  const handleJoin = (e: React.FormEvent) => {
    e.preventDefault();
    if (username.trim()) {
        joinGame(username)
    }
  };

  if (user) {
    return null
  }

  return (
    <div className="absolute inset-0 bg-black/75 flex items-center justify-center z-10">
      <form onSubmit={handleJoin} className="bg-gray-800 p-8 rounded-lg shadow-lg flex flex-col gap-4">
        <h2 className="text-2xl font-bold text-white">Enter the Arena</h2>
        <input
          type="text"
          value={username}
          onChange={(e) => setUsername(e.target.value)}
          placeholder="Enter your name"
          className="p-2 rounded-md bg-gray-700 text-white border border-gray-600 focus:outline-none focus:ring-2 focus:ring-blue-500"
        />
        <button type="submit" className="bg-blue-600 text-white p-2 rounded-md hover:bg-blue-700 transition-colors">
          Join Game
        </button>
      </form>
    </div>
  );
};

export default JoinGameOverlay;
