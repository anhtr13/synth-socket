import { z } from "zod";

const LoginSchema = z.object({
  user_email: z.email(),
  password: z.string().min(6, {
    message: "Password must be 6 characters or longer",
  }),
});
type LoginPayload = z.infer<typeof LoginSchema>;

const SignupSchema = z.object({
  user_email: z.email(),
  user_name: z
    .string()
    .min(3, {
      message: "User name must be 3 characters or longer",
    })
    .max(30, {
      message: "User name must be less than 30 characters",
    }),
  password: z.string().min(6, {
    message: "Password must be 6 characters or longer",
  }),
});
type SignupPayload = z.infer<typeof SignupSchema>;
export { LoginSchema, SignupSchema, type LoginPayload, type SignupPayload };
