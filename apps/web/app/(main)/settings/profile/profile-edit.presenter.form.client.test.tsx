import { expect, test as base, vi } from "vitest";
import { render } from "vitest-browser-react";
import type { UpdateProfileResult } from "./profile-edit.action";
import { ProfileEditPresenter } from "./profile-edit.presenter.form.client";

const editProfileMock =
  vi.fn<(nickname: string, introduce: string) => Promise<UpdateProfileResult>>();

vi.mock("./profile-edit.action", () => ({
  editProfile: (nickname: string, introduce: string) =>
    editProfileMock(nickname, introduce),
}));

const test = base.extend({
  editProfile: [
    async ({}, provide: (v: typeof editProfileMock) => Promise<void>) => {
      editProfileMock.mockReset();
      await provide(editProfileMock);
    },
    { auto: true },
  ],
});

test("初期値が入力欄に表示される", async () => {
  const screen = await render(
    <ProfileEditPresenter defaultNickname="taro" defaultIntroduce="hello" />,
  );

  await expect.element(screen.getByLabelText(/ニックネーム/)).toHaveValue("taro");
  await expect.element(screen.getByLabelText(/自己紹介/)).toHaveValue("hello");
});

test("空のニックネームで blur すると必須エラーが表示される", async () => {
  const screen = await render(<ProfileEditPresenter defaultNickname="taro" />);

  const nickname = screen.getByLabelText(/ニックネーム/);
  await nickname.clear();
  await screen.getByLabelText(/自己紹介/).click();

  await expect
    .element(screen.getByText("ニックネームは必須です"))
    .toBeInTheDocument();
});

test("ニックネーム 50 文字（上限）は受理される", async () => {
  const screen = await render(<ProfileEditPresenter defaultNickname="taro" />);

  const nickname = screen.getByLabelText(/ニックネーム/);
  await nickname.fill("a".repeat(50));
  await screen.getByLabelText(/自己紹介/).click();

  await expect
    .element(screen.getByText("ニックネームは50文字以内で入力してください"))
    .not.toBeInTheDocument();
});

test("ニックネーム 51 文字で blur すると文字数超過エラーが表示される", async () => {
  const screen = await render(<ProfileEditPresenter defaultNickname="taro" />);

  const nickname = screen.getByLabelText(/ニックネーム/);
  await nickname.fill("a".repeat(51));
  await screen.getByLabelText(/自己紹介/).click();

  await expect
    .element(screen.getByText("ニックネームは50文字以内で入力してください"))
    .toBeInTheDocument();
});

test("自己紹介 500 文字（上限）は受理される", async () => {
  const screen = await render(<ProfileEditPresenter defaultNickname="taro" />);

  const introduce = screen.getByLabelText(/自己紹介/);
  await introduce.fill("a".repeat(500));
  await screen.getByLabelText(/ニックネーム/).click();

  await expect
    .element(screen.getByText("自己紹介は500文字以内で入力してください"))
    .not.toBeInTheDocument();
});

test("自己紹介 501 文字で blur すると文字数超過エラーが表示される", async () => {
  const screen = await render(<ProfileEditPresenter defaultNickname="taro" />);

  const introduce = screen.getByLabelText(/自己紹介/);
  await introduce.fill("a".repeat(501));
  await screen.getByLabelText(/ニックネーム/).click();

  await expect
    .element(screen.getByText("自己紹介は500文字以内で入力してください"))
    .toBeInTheDocument();
});

test("送信成功時に成功メッセージが表示され、入力値で action が呼ばれる", async () => {
  editProfileMock.mockResolvedValue({ success: true });
  const screen = await render(
    <ProfileEditPresenter defaultNickname="taro" defaultIntroduce="hi" />,
  );

  await screen.getByRole("button", { name: "保存" }).click();

  await expect
    .element(screen.getByText("プロフィールを保存しました"))
    .toBeInTheDocument();
  expect(editProfileMock).toHaveBeenCalledWith("taro", "hi");
});

test("UNAUTHORIZED コードを返した場合はログイン要求メッセージが表示される", async () => {
  editProfileMock.mockResolvedValue({ success: false, code: "UNAUTHORIZED" });
  const screen = await render(<ProfileEditPresenter defaultNickname="taro" />);

  await screen.getByRole("button", { name: "保存" }).click();

  await expect
    .element(screen.getByText("ログインが必要です"))
    .toBeInTheDocument();
  await expect
    .element(screen.getByText("プロフィールを保存しました"))
    .not.toBeInTheDocument();
});

test("UPDATE_FAILED コードを返した場合は更新失敗メッセージが表示される", async () => {
  editProfileMock.mockResolvedValue({ success: false, code: "UPDATE_FAILED" });
  const screen = await render(<ProfileEditPresenter defaultNickname="taro" />);

  await screen.getByRole("button", { name: "保存" }).click();

  await expect
    .element(screen.getByText("プロフィールの更新に失敗しました"))
    .toBeInTheDocument();
});

test("action が例外を投げた場合は汎用の失敗メッセージが表示される", async () => {
  editProfileMock.mockRejectedValueOnce(new Error("network down"));
  const screen = await render(<ProfileEditPresenter defaultNickname="taro" />);

  await screen.getByRole("button", { name: "保存" }).click();

  await expect
    .element(screen.getByText("プロフィールの更新に失敗しました"))
    .toBeInTheDocument();
  await expect
    .element(screen.getByText("プロフィールを保存しました"))
    .not.toBeInTheDocument();
  expect(editProfileMock).toHaveBeenCalledTimes(1);
});

test("バリデーションエラーがある状態で保存ボタンを押しても action は呼ばれない", async () => {
  const screen = await render(<ProfileEditPresenter defaultNickname="taro" />);

  const nickname = screen.getByLabelText(/ニックネーム/);
  await nickname.clear();
  await screen.getByLabelText(/自己紹介/).click();

  await expect
    .element(screen.getByText("ニックネームは必須です"))
    .toBeInTheDocument();

  // 押下を阻止する仕組み（aria-disabled / 内部バリデーション 等）に依存せず、
  // 「保存しようとしても action が呼ばれない」ことだけを検証する
  await screen
    .getByRole("button", { name: "保存" })
    .click({ force: true });

  await new Promise((resolve) => setTimeout(resolve, 50));
  expect(editProfileMock).not.toHaveBeenCalled();
});

test("disabled=true のときは入力欄と送信ボタンが無効化される", async () => {
  const screen = await render(<ProfileEditPresenter disabled />);

  await expect.element(screen.getByLabelText(/ニックネーム/)).toBeDisabled();
  await expect.element(screen.getByLabelText(/自己紹介/)).toBeDisabled();
  await expect.element(screen.getByRole("button", { name: "保存" })).toBeDisabled();
});
