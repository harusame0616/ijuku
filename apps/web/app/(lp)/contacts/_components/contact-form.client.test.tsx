import { expect, test as base, vi } from "vitest";
import { render } from "vitest-browser-react";
import type { SubmitContactResult } from "@/app/(lp)/contacts/contact.action";
import { ContactForm } from "@/app/(lp)/contacts/_components/contact-form.client";

const submitContactMock = vi.fn<
  (input: {
    name: string;
    email: string;
    phone: string;
    content: string;
  }) => Promise<SubmitContactResult>
>();

vi.mock("@/app/(lp)/contacts/contact.action", () => ({
  submitContact: (input: {
    name: string;
    email: string;
    phone: string;
    content: string;
  }) => submitContactMock(input),
}));

const test = base.extend({
  submitContact: [
    async ({}, provide: (v: typeof submitContactMock) => Promise<void>) => {
      submitContactMock.mockReset();
      await provide(submitContactMock);
    },
    { auto: true },
  ],
});

async function fillRequiredFields(screen: Awaited<ReturnType<typeof render>>) {
  await screen.getByLabelText(/お名前/).fill("山田 太郎");
  await screen.getByLabelText(/メールアドレス/).fill("taro@example.com");
  await screen.getByLabelText(/お問い合わせ内容/).fill("質問があります。");
  // 入力欄から focus を外して blur を確定させ、 onBlur バリデーションを走らせる
  await screen.getByLabelText(/お名前/).click();
  await screen.getByLabelText(/メールアドレス/).click();
}

test("必須項目のラベルに「必須」、任意項目に「任意」が表示される", async () => {
  const screen = await render(<ContactForm />);

  await expect
    .element(screen.getByLabelText(/お名前\s*必須/))
    .toBeInTheDocument();
  await expect
    .element(screen.getByLabelText(/メールアドレス\s*必須/))
    .toBeInTheDocument();
  await expect
    .element(screen.getByLabelText(/電話番号\s*任意/))
    .toBeInTheDocument();
  await expect
    .element(screen.getByLabelText(/お問い合わせ内容\s*必須/))
    .toBeInTheDocument();
});

test("空のお名前で blur すると必須エラーが表示される", async () => {
  const screen = await render(<ContactForm />);

  const name = screen.getByLabelText(/お名前/);
  await name.click();
  await screen.getByLabelText(/メールアドレス/).click();

  await expect
    .element(screen.getByText("お名前は必須です"))
    .toBeInTheDocument();
});

test("不正なメールアドレスを入れて blur すると形式エラーが表示される", async () => {
  const screen = await render(<ContactForm />);

  await screen.getByLabelText(/メールアドレス/).fill("not-email");
  await screen.getByLabelText(/お名前/).click();

  await expect
    .element(screen.getByText("メールアドレスの形式が正しくありません"))
    .toBeInTheDocument();
});

test("送信成功時に完了メッセージが表示され、入力値で action が呼ばれる", async () => {
  submitContactMock.mockResolvedValue({ success: true });
  const screen = await render(<ContactForm />);

  await fillRequiredFields(screen);
  await screen.getByLabelText(/電話番号/).fill("090-1234-5678");
  await screen.getByRole("button", { name: "送信する" }).click();

  await expect
    .element(screen.getByText("お問い合わせを受け付けました。"))
    .toBeInTheDocument();
  expect(submitContactMock).toHaveBeenCalledWith({
    name: "山田 太郎",
    email: "taro@example.com",
    phone: "090-1234-5678",
    content: "質問があります。",
  });
});

test("送信成功後はフォームの入力値がリセットされる", async () => {
  submitContactMock.mockResolvedValue({ success: true });
  const screen = await render(<ContactForm />);

  await fillRequiredFields(screen);
  await screen.getByRole("button", { name: "送信する" }).click();

  await expect
    .element(screen.getByText("お問い合わせを受け付けました。"))
    .toBeInTheDocument();

  await expect.element(screen.getByLabelText(/お名前/)).toHaveValue("");
  await expect.element(screen.getByLabelText(/メールアドレス/)).toHaveValue("");
  await expect
    .element(screen.getByLabelText(/お問い合わせ内容/))
    .toHaveValue("");
});

test("SUBMIT_FAILED コードを返した場合は失敗メッセージが表示される", async () => {
  submitContactMock.mockResolvedValue({ success: false, code: "SUBMIT_FAILED" });
  const screen = await render(<ContactForm />);

  await fillRequiredFields(screen);
  await screen.getByRole("button", { name: "送信する" }).click();

  await expect
    .element(
      screen.getByText("送信に失敗しました。時間をおいて再度お試しください。"),
    )
    .toBeInTheDocument();
  await expect
    .element(screen.getByText("お問い合わせを受け付けました。"))
    .not.toBeInTheDocument();
});

test("action が例外を投げた場合は失敗メッセージが表示される", async () => {
  submitContactMock.mockRejectedValueOnce(new Error("network down"));
  const screen = await render(<ContactForm />);

  await fillRequiredFields(screen);
  await screen.getByRole("button", { name: "送信する" }).click();

  await expect
    .element(
      screen.getByText("送信に失敗しました。時間をおいて再度お試しください。"),
    )
    .toBeInTheDocument();
});

test("バリデーションエラーがある状態で送信ボタンを押しても action は呼ばれない", async () => {
  const screen = await render(<ContactForm />);

  // 空のまま blur で必須エラーを発生させる
  await screen.getByLabelText(/お名前/).click();
  await screen.getByLabelText(/メールアドレス/).click();

  await screen
    .getByRole("button", { name: "送信する" })
    .click({ force: true });

  await new Promise((resolve) => setTimeout(resolve, 50));
  expect(submitContactMock).not.toHaveBeenCalled();
});
