import React from 'react';

const Leaderboard = () => {
  const players = [
    { id: 1, name: 'Player 1', score: 150000 },
    { id: 2, name: 'Player 2', score: 145000 },
    { id: 3, name: 'Player 3', score: 130000 },
    { id: 4, name: 'Player 4', score: 125000 },
    { id: 5, name: 'Player 5', score: 120000 },
  ];

  return (
    <div className="p-4 rounded-lg">
      <h2 className="text-lg font-semibold text-gray-300 mb-2">Leaderboard</h2>
      <ul>
        {players.map((player) => (
          <li key={player.id} className="flex justify-between items-center text-sm text-gray-400 py-1">
            <span>{player.name}</span>
            <span>${player.score.toLocaleString()}</span>
          </li>
        ))}
      </ul>
    </div>
  );
};

export default Leaderboard; 