"use client";

import apiService from "@/api";
import useAuthStore from "@/stores/useAuthStore";
import { useWsStore } from "@/stores/useWsStore";
import { useEffect } from "react";

const Provider = ({ children }: { children: React.ReactNode }) => {
  const { user, token, isTokenExpired, refreshAuthToken, refreshToken } =
    useAuthStore();
  const { connect, disconnect } = useWsStore();

  useEffect(() => {
    if (token) {
      apiService.setToken(token);
    }
  }, [token]);

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
  }, [user, token, connect, disconnect]);

  return <>{children}</>;
};

export default Provider;
