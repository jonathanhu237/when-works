import { createRootRouteWithContext, Outlet } from "@tanstack/react-router";
import { TanStackRouterDevtools } from "@tanstack/react-router-devtools";
import { ReactQueryDevtools } from "@tanstack/react-query-devtools";
import type { QueryClient } from "@tanstack/react-query";
import { Toaster } from "@/components/ui/sonner";

const RootLayout = () => (
  <>
    <Outlet />
    <TanStackRouterDevtools />
    <ReactQueryDevtools />
    <Toaster />
  </>
);

export const Route = createRootRouteWithContext<{
  queryClient: QueryClient;
}>()({ component: RootLayout });
