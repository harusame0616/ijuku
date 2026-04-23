import { Badge } from "@/components/ui/badge";

type Props = {
  required: boolean;
};

export function RequiredOptionalBadge({ required }: Props) {
  return required ? (
    <Badge variant="outline" className="border-primary text-primary">
      必須
    </Badge>
  ) : (
    <Badge variant="outline" className="border-secondary text-secondary">
      任意
    </Badge>
  );
}
