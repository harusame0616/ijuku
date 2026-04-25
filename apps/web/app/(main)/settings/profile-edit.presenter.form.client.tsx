"use client";

import { Button } from "@/components/ui/button";
import {
  Field,
  FieldDescription,
  FieldError,
  FieldLabel,
} from "@/components/ui/field";
import { Input } from "@/components/ui/input";
import { RequiredOptionalBadge } from "@/components/ui/required-optional-badge";
import { Textarea } from "@/components/ui/textarea";
import { useIsHydrated } from "@/hooks/use-is-hydrated";
import { useForm } from "@tanstack/react-form";
import { useState } from "react";
import * as v from "valibot";
import { editProfile, type EditProfileErrorCode } from "./profile-edit.action";

const errorMessages: Record<EditProfileErrorCode, string> = {
  UNAUTHORIZED: "ログインが必要です",
  UPDATE_FAILED: "プロフィールの更新に失敗しました",
};

const nicknameSchema = v.pipe(
  v.string(),
  v.minLength(1, "ニックネームは必須です"),
  v.maxLength(50, "ニックネームは50文字以内で入力してください"),
);

const introduceSchema = v.pipe(
  v.string(),
  v.maxLength(500, "自己紹介は500文字以内で入力してください"),
);

const schema = v.object({
  nickname: nicknameSchema,
  introduce: introduceSchema,
});

type Props = {
  defaultNickname?: string;
  defaultIntroduce?: string;
  disabled?: boolean;
};

export function ProfileEditPresenter({
  defaultNickname = "",
  defaultIntroduce = "",
  disabled: disabledProp = false,
}: Props) {
  const isHydrated = useIsHydrated();
  const disabled = disabledProp || !isHydrated;
  const [successMessage, setSuccessMessage] = useState<string | null>(null);
  const [submitError, setSubmitError] = useState<string | null>(null);

  const form = useForm({
    defaultValues: {
      nickname: defaultNickname,
      introduce: defaultIntroduce,
    },
    validators: {
      onBlur: schema,
    },
    onSubmit: async ({ value }) => {
      setSuccessMessage(null);
      setSubmitError(null);
      try {
        const res = await editProfile(value.nickname, value.introduce);
        if (res.success) {
          setSuccessMessage("プロフィールを保存しました");
        } else {
          setSubmitError(errorMessages[res.code]);
        }
      } catch {
        setSubmitError(errorMessages.UPDATE_FAILED);
      }
    },
  });

  return (
    <form
      onSubmit={(e) => {
        e.preventDefault();
        e.stopPropagation();
        form.handleSubmit();
      }}
      className="grid gap-6"
    >
      <form.Field name="nickname">
        {(field) => {
          const isInvalid = !field.state.meta.isValid;
          return (
            <Field data-invalid={isInvalid}>
              <FieldLabel htmlFor={field.name}>
                ニックネーム
                <RequiredOptionalBadge required />
              </FieldLabel>
              <Input
                id={field.name}
                name={field.name}
                aria-invalid={isInvalid}
                value={field.state.value}
                onChange={(e) => field.handleChange(e.target.value)}
                onBlur={field.handleBlur}
                disabled={disabled}
              />
              <FieldDescription>50文字以内</FieldDescription>
              <FieldError errors={field.state.meta.errors} />
            </Field>
          );
        }}
      </form.Field>

      <form.Field name="introduce">
        {(field) => {
          const isInvalid = !field.state.meta.isValid;
          return (
            <Field data-invalid={isInvalid}>
              <FieldLabel htmlFor={field.name}>
                自己紹介
                <RequiredOptionalBadge required={false} />
              </FieldLabel>
              <Textarea
                id={field.name}
                name={field.name}
                aria-invalid={isInvalid}
                value={field.state.value}
                onChange={(e) => field.handleChange(e.target.value)}
                onBlur={field.handleBlur}
                className="h-40"
                disabled={disabled}
              />
              <FieldDescription>500文字以内</FieldDescription>
              <FieldError errors={field.state.meta.errors} />
            </Field>
          );
        }}
      </form.Field>

      {submitError && <p className="text-sm text-destructive">{submitError}</p>}

      {successMessage && (
        <p className="text-sm text-primary">{successMessage}</p>
      )}

      <form.Subscribe selector={(state) => [state.isSubmitting, state.isValid]}>
        {([isSubmitting, isValid]) => (
          <Button
            type="submit"
            disabled={disabled || isSubmitting || !isValid}
            className="mt-8"
          >
            {isSubmitting ? "保存中..." : "保存"}
          </Button>
        )}
      </form.Subscribe>
    </form>
  );
}
