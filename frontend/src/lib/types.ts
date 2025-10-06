export interface ApiError<T = unknown> {
  code: string;
  message: string;
  details: T | null;
}

export function isApiError(error: unknown): error is ApiError {
  return (
    typeof error === "object" &&
    error !== null &&
    "code" in error &&
    "message" in error &&
    "details" in error
  );
}

export interface User {
  id: string;
  username: string;
  email: string;
  name: string;
  is_admin: boolean;
  created_at: string;
}
