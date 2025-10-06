import axios, { AxiosError } from "axios";
import { isApiError, type ApiError, type User } from "./types";
import { queryOptions } from "@tanstack/react-query";

const api = axios.create({
  baseURL: "/api",
  withCredentials: true,
  headers: {
    "Content-Type": "application/json",
  },
});

api.interceptors.response.use(
  (response) => response,
  (error: AxiosError) => {
    if (error.response && isApiError(error.response.data)) {
      return Promise.reject(error.response.data);
    }

    const unknownError: ApiError = {
      code: "UNKNOWN_ERROR",
      message: "An unknown error occurred",
      details: null,
    };
    return Promise.reject(unknownError);
  }
);

// --------------------------------------------------
// Get my profile
// --------------------------------------------------
export async function getMyProfile(): Promise<User | null> {
  try {
    const res = await api.get<{ user: User }>("/v1/me");
    return res.data.user;
  } catch (error) {
    if (isApiError(error) && error.code === "UNAUTHORIZED") {
      return null;
    }
    throw error;
  }
}

export const getMyProfileQueryOptions = () =>
  queryOptions({
    queryKey: ["me"],
    queryFn: getMyProfile,
    retry: false,
  });

// --------------------------------------------------
// Login
// --------------------------------------------------
export async function login(data: {
  username: string;
  password: string;
}): Promise<User> {
  const res = await api.post<{ user: User }>("/v1/auth/login", data);
  return res.data.user;
}
