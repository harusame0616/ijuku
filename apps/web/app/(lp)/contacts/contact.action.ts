"use server";

import { headers } from "next/headers";

export type SubmitContactErrorCode = "SUBMIT_FAILED";

export type SubmitContactResult =
  | { success: true }
  | { success: false; code: SubmitContactErrorCode };

interface ContactInput {
  name: string;
  email: string;
  phone: string;
  content: string;
}

export async function submitContact(
  input: ContactInput,
): Promise<SubmitContactResult> {
  const requestHeaders = await headers();
  const forwardedFor = requestHeaders.get("x-forwarded-for");
  const userAgent = requestHeaders.get("user-agent") ?? "";

  console.log(Object.fromEntries(requestHeaders.entries()))
  const phone = input.phone.trim();

  const response = await fetch(`${process.env.API_URL}/v1/contacts`, {
    method: "POST",
    headers: {
      "Content-Type": "application/json",
      ...(forwardedFor ? { "X-Forwarded-For": forwardedFor } : {}),
      "User-Agent": userAgent,
    },
    body: JSON.stringify({
      name: input.name,
      email: input.email,
      phone: phone === "" ? null : phone,
      content: input.content,
    }),
  });

  if (!response.ok) {
    return { success: false, code: "SUBMIT_FAILED" };
  }

  return { success: true };
}
