import { getMyProfileQueryOptions } from "@/lib/api";
import { createFileRoute, Outlet, redirect } from "@tanstack/react-router";

export const Route = createFileRoute("/_unauthenticated")({
  beforeLoad: async ({ context }) => {
    const user = await context.queryClient.ensureQueryData(
      getMyProfileQueryOptions()
    );

    if (user) {
      throw redirect({ to: "/" });
    }
  },
  component: RouteComponent,
});

function RouteComponent() {
  return <Outlet />;
}
