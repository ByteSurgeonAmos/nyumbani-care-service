import { error } from "@sveltejs/kit";

export const load = ({ params }: { params: { id: string } }) => {
  // Get the test kit ID from the URL parameter
  const { id } = params;

  if (!id) {
    throw error(404, {
      message: "Test kit not found",
    });
  }

  return {
    id,
  };
};
