"use client";

import type { JSX } from "react";
import { useState } from "react";
import { useForm } from "@tanstack/react-form";
import * as v from "valibot";
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
import {
  submitContact,
  type SubmitContactErrorCode,
} from "@/app/(lp)/contacts/contact.action";

const errorMessages: Record<SubmitContactErrorCode, string> = {
  SUBMIT_FAILED: "送信に失敗しました。時間をおいて再度お試しください。",
};

const nameSchema = v.pipe(
  v.string(),
  v.minLength(1, "お名前は必須です"),
  v.maxLength(100, "お名前は100文字以内で入力してください"),
);

const emailSchema = v.pipe(
  v.string(),
  v.minLength(1, "メールアドレスは必須です"),
  v.maxLength(255, "メールアドレスは255文字以内で入力してください"),
  v.email("メールアドレスの形式が正しくありません"),
);

const phoneSchema = v.pipe(
  v.string(),
  v.maxLength(20, "電話番号は20文字以内で入力してください"),
);

const contentSchema = v.pipe(
  v.string(),
  v.minLength(1, "お問い合わせ内容は必須です"),
  v.maxLength(2000, "お問い合わせ内容は2000文字以内で入力してください"),
);

const schema = v.object({
  name: nameSchema,
  email: emailSchema,
  phone: phoneSchema,
  content: contentSchema,
});

const defaultValues = {
  name: "",
  email: "",
  phone: "",
  content: "",
};

export function ContactForm(): JSX.Element {
  const isHydrated = useIsHydrated();
  const [successMessage, setSuccessMessage] = useState<string | null>(null);
  const [submitError, setSubmitError] = useState<string | null>(null);

  const form = useForm({
    defaultValues,
    validators: {
      onBlur: schema,
    },
    onSubmit: async ({ value, formApi }) => {
      setSuccessMessage(null);
      setSubmitError(null);
      try {
        const response = await submitContact(value);
        if (response.success) {
          setSuccessMessage("お問い合わせを受け付けました。");
          formApi.reset();
        } else {
          setSubmitError(errorMessages[response.code]);
        }
      } catch {
        setSubmitError(errorMessages.SUBMIT_FAILED);
      }
    },
  });

  return (
    <form
      onSubmit={(event) => {
        event.preventDefault();
        event.stopPropagation();
        form.handleSubmit();
      }}
      className="grid gap-6"
      noValidate
    >
      <form.Field name="name">
        {(field) => {
          const isInvalid = !field.state.meta.isValid;
          return (
            <Field data-invalid={isInvalid}>
              <FieldLabel htmlFor={field.name}>
                お名前
                <RequiredOptionalBadge required />
              </FieldLabel>
              <Input
                id={field.name}
                name={field.name}
                aria-invalid={isInvalid}
                value={field.state.value}
                onChange={(event) => field.handleChange(event.target.value)}
                onBlur={field.handleBlur}
                disabled={!isHydrated}
              />
              <FieldDescription>100文字以内</FieldDescription>
              <FieldError errors={field.state.meta.errors} />
            </Field>
          );
        }}
      </form.Field>

      <form.Field name="email">
        {(field) => {
          const isInvalid = !field.state.meta.isValid;
          return (
            <Field data-invalid={isInvalid}>
              <FieldLabel htmlFor={field.name}>
                メールアドレス
                <RequiredOptionalBadge required />
              </FieldLabel>
              <Input
                id={field.name}
                name={field.name}
                type="email"
                inputMode="email"
                autoComplete="email"
                aria-invalid={isInvalid}
                value={field.state.value}
                onChange={(event) => field.handleChange(event.target.value)}
                onBlur={field.handleBlur}
                disabled={!isHydrated}
              />
              <FieldDescription>
                返信に使用します。255文字以内。
              </FieldDescription>
              <FieldError errors={field.state.meta.errors} />
            </Field>
          );
        }}
      </form.Field>

      <form.Field name="phone">
        {(field) => {
          const isInvalid = !field.state.meta.isValid;
          return (
            <Field data-invalid={isInvalid}>
              <FieldLabel htmlFor={field.name}>
                電話番号
                <RequiredOptionalBadge required={false} />
              </FieldLabel>
              <Input
                id={field.name}
                name={field.name}
                type="tel"
                inputMode="tel"
                autoComplete="tel"
                aria-invalid={isInvalid}
                value={field.state.value}
                onChange={(event) => field.handleChange(event.target.value)}
                onBlur={field.handleBlur}
                disabled={!isHydrated}
              />
              <FieldDescription>20文字以内</FieldDescription>
              <FieldError errors={field.state.meta.errors} />
            </Field>
          );
        }}
      </form.Field>

      <form.Field name="content">
        {(field) => {
          const isInvalid = !field.state.meta.isValid;
          return (
            <Field data-invalid={isInvalid}>
              <FieldLabel htmlFor={field.name}>
                お問い合わせ内容
                <RequiredOptionalBadge required />
              </FieldLabel>
              <Textarea
                id={field.name}
                name={field.name}
                aria-invalid={isInvalid}
                value={field.state.value}
                onChange={(event) => field.handleChange(event.target.value)}
                onBlur={field.handleBlur}
                className="h-48"
                disabled={!isHydrated}
              />
              <FieldDescription>2000文字以内</FieldDescription>
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
            disabled={!isHydrated || isSubmitting || !isValid}
            className="mt-4"
          >
            {isSubmitting ? "送信中..." : "送信する"}
          </Button>
        )}
      </form.Subscribe>
    </form>
  );
}
