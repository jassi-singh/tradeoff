"use client";
import useAuthStore from "@/stores/useAuthStore";
import React, { useEffect } from "react";

const JoinGameOverlay: React.FC<{ children: React.ReactNode }> = ({
  children,
}) => {
  const [isHydrated, setIsHydrated] = React.useState(false);
  const { user, joinGame } = useAuthStore();

  useEffect(() => {
    setIsHydrated(true);
  }, []);

  const handleJoin = (e: React.FormEvent) => {
    e.preventDefault();
    const formData = new FormData(e.target as HTMLFormElement);
    const username = formData.get("username") as string;
    if (username.trim()) {
      joinGame(username);
    }
  };

  if (user || !isHydrated) {
    return null;
  }

  return (
    <div className="absolute inset-0 bg-black/60 backdrop-blur-md flex items-center justify-center z-10">
      <div className="relative">
        <form
          onSubmit={handleJoin}
          className="bg-gray-900/95 backdrop-blur-sm p-8 rounded-xl shadow-2xl border border-gray-700/50 max-w-md w-full mx-4"
        >
          {children}
        </form>
      </div>
    </div>
  );
};

export default JoinGameOverlay;
