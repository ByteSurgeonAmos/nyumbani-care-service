import { error } from "@sveltejs/kit";

export const load = ({ params }: { params: { id: string } }) => {
  const { id } = params;

  if (!id) {
    throw error(404, {
      message: "Article not found",
    });
  }

  return {
    id,
  };
};
