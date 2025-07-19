"use client";

import useAuthStore from "@/stores/useAuthStore";
import { useWsStore } from "@/stores/useWsStore";
import { useEffect } from "react";

const Provider = ({ children }: { children: React.ReactNode }) => {
  const { user, token } = useAuthStore();
  const { connect, disconnect } = useWsStore();

  useEffect(() => {
    if (user && token) {
      connect();
    } else {
      disconnect();
    }

    return () => {
      disconnect();
    };
  }, [user, token, connect, disconnect]);

  return <>{children}</>;
};

export default Provider;
