import { createUser } from "@/api";
import { User } from "@/types";
import { create } from "zustand";
import { persist } from "zustand/middleware";

interface AuthStore {
    user: User | null
    joinGame: (username: string) => Promise<void>;
}

const useAuthStore = create<AuthStore>()(
    persist(
        (set) => ({
            user: null,
            joinGame: async (username: string) => {
                const user = await createUser(username)
                set({ user })
            }
        }),
        { name: "current-user-storage" }
    )
)

export default useAuthStore;