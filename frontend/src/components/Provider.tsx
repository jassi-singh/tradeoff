"use client";

import useAuthStore from "@/stores/useAuthStore";
import { useWsStore } from "@/stores/useWsStore";
import { useEffect } from "react";

const Provider = ({ children }: { children: React.ReactNode }) => {
  const {
    user,
    token,
    isTokenExpired,
    refreshAuthToken,
    refreshToken,
    initializeAuth,
  } = useAuthStore();
  const { connect, disconnect } = useWsStore();

  // Initialize auth on mount
  useEffect(() => {
    initializeAuth();
  }, [initializeAuth]);

  useEffect(() => {
    if (isTokenExpired() && refreshToken) {
      refreshAuthToken();
    }
  }, [isTokenExpired, refreshAuthToken, refreshToken]);

  useEffect(() => {
    if (user && token && !isTokenExpired()) {
      connect();
    } else {
      disconnect();
    }

    return () => {
      disconnect();
    };
  }, [user, token, connect, disconnect, isTokenExpired]);

  return <>{children}</>;
};

export default Provider;
