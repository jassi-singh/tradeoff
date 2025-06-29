"use client";
import useAuthStore from "@/stores/useAuthStore";
import React from "react";

const JoinGameOverlay: React.FC<{ children: React.ReactNode }> = ({
  children,
}) => {
  const { user, joinGame } = useAuthStore();

  const handleJoin = (e: React.FormEvent) => {
    e.preventDefault();
    const formData = new FormData(e.target as HTMLFormElement);
    const username = formData.get("username") as string;
    if (username.trim()) {
      joinGame(username);
    }
  };

  if (user) {
    return null;
  }

  return (
    <div className="absolute inset-0 bg-black/75 flex items-center justify-center z-10">
      <form
        onSubmit={handleJoin}
        className="bg-gray-800 p-8 rounded-lg shadow-lg flex flex-col gap-4"
      >
        {children}
      </form>
    </div>
  );
};

export default JoinGameOverlay;
