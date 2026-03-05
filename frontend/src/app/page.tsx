import { redirect } from "next/navigation";

export default function Home() {
  // In a real app we might check for an auth cookie here
  redirect("/feed");
}
