"use client";

import useAuthStore from "@/stores/useAuthStore";
import { useWsStore } from "@/stores/useWsStore";
import { useEffect } from "react";

const Provider = ({ children }: { children: React.ReactNode }) => {
  const { user } = useAuthStore();
  const { connect, disconnect } = useWsStore();

  useEffect(() => {
    if (user) {
      connect(user.id);
    }

    return () => {
      disconnect();
    };
  }, [user]);

  return <>{children}</>;
};

export default Provider;
