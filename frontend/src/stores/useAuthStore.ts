import { login } from "@/api";
import { User } from "@/types";
import { create } from "zustand";
import { persist } from "zustand/middleware";

interface AuthStore {
    user: User | null
    token: string | null;
    refreshToken: string | null;
    joinGame: (username: string) => Promise<void>;
}

const useAuthStore = create<AuthStore>()(
    persist(
        (set) => ({
            user: null,
            token: null,
            refreshToken: null,
            joinGame: async (username: string) => {
                try {
                    const response = await login(username)
                    set({ user: response.user, token: response.token, refreshToken: response.refreshToken });
                } catch (error) {
                    console.error("Failed to join game:", error);
                    set({ user: null });
                }
            }
        }),
        { name: "current-user-storage" }
    )
)

export default useAuthStore;
