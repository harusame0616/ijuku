import type { JSX } from "react";
import type { Metadata } from "next";

export const metadata: Metadata = {
  title: "プライバシーポリシー | JukuBox.ai",
  description:
    "JukuBox.ai のプライバシーポリシー。取得する個人情報等、利用目的、第三者提供、安全管理措置、保存期間、お問い合わせ窓口について記載しています。",
};

export default function PrivacyPoliciesPage(): JSX.Element {
  return (
    <main className="bg-background min-h-screen">
      <article className="mx-auto max-w-3xl px-6 py-16 text-foreground md:px-8 md:py-24">
        <h1 className="font-serif text-3xl md:text-4xl">プライバシーポリシー</h1>
        <p className="mt-4 leading-relaxed text-muted-foreground">
          はるさめ dev（以下「当方」といいます。）は、当方が運営する AI 学習プラットフォーム
          JukuBox.ai（以下「本サービス」といいます。）におけるお客様の個人情報等の取扱いについて、以下のとおりプライバシーポリシー（以下「本ポリシー」といいます。）を定めます。本ポリシーにおいて「個人情報等」とは、個人情報の保護に関する法律（以下「個人情報保護法」といいます。）に定める個人情報及び個人関連情報を総称していいます。
        </p>

        <h2 className="mt-12 border-b border-border pb-2 font-serif text-xl md:text-2xl">
          1. 取得する個人情報等
        </h2>
        <p className="mt-4 leading-relaxed text-muted-foreground">
          当方が取得する個人情報等は、その取得方法に応じて、以下のようなものがあります。
        </p>

        <h3 className="mt-6 text-base font-semibold">
          (1) お客様から直接ご提供いただく情報
        </h3>
        <p className="mt-2 leading-relaxed text-muted-foreground">
          お客様が本サービス上の入力フォーム等を通じて、当方に直接ご提供いただく情報です。
        </p>
        <ul className="mt-3 list-disc space-y-1 pl-6 text-muted-foreground">
          <li>ニックネーム、自己紹介文その他プロフィールに関する情報</li>
          <li>JukuBox API キー（お客様の操作により生成された情報）</li>
          <li>お問い合わせフォームを通じて送信いただく情報</li>
        </ul>

        <h3 className="mt-6 text-base font-semibold">
          (2) 第三者からご提供いただく情報
        </h3>
        <p className="mt-2 leading-relaxed text-muted-foreground">
          お客様が外部サービスとの連携を許可した場合に、当該外部サービスから当方が取得する情報です。
        </p>
        <ul className="mt-3 list-disc space-y-1 pl-6 text-muted-foreground">
          <li>Google アカウントによる認証時に取得するメールアドレス等のプロフィール情報</li>
        </ul>

        <h3 className="mt-6 text-base font-semibold">
          (3) お客様が本サービスを利用するにあたって、当方が自動的に取得する情報
        </h3>
        <p className="mt-2 leading-relaxed text-muted-foreground">
          お客様が本サービスを利用する際に、システムが自動的に取得・記録する情報です。
        </p>
        <ul className="mt-3 list-disc space-y-1 pl-6 text-muted-foreground">
          <li>セッション維持のための Cookie に関する情報</li>
          <li>最終ログイン日時</li>
          <li>本サービスのご利用状況に関する情報（コース受講履歴、トピック学習進捗等）</li>
          <li>
            IP アドレス、ブラウザの種類及びバージョン（User-Agent）、リファラー、アクセス日時、閲覧 URL 等のアクセスログ
          </li>
          <li>
            Google Analytics（GA4）が発行する Cookie 識別子（Client ID）等のアクセス解析に関する情報
          </li>
        </ul>

        <h3 className="mt-6 text-base font-semibold">
          (4) お客様が著者として登録された場合に取得する情報（公表前提）
        </h3>
        <p className="mt-2 leading-relaxed text-muted-foreground">
          お客様が本サービス上で著者として登録された場合に取得する、公表を目的とした情報です。これらの情報は、本サービス上で広く一般に公開され、本サービスを利用する全てのお客様が閲覧できます。お客様は、これらが公開情報であることに同意したうえで登録するものとします。
        </p>
        <ul className="mt-3 list-disc space-y-1 pl-6 text-muted-foreground">
          <li>著者名</li>
          <li>著者プロフィール</li>
          <li>著者スラッグ（著者ページの URL に使用される識別子）</li>
        </ul>

        <h2 className="mt-12 border-b border-border pb-2 font-serif text-xl md:text-2xl">
          2. 個人情報等の利用目的
        </h2>
        <p className="mt-4 leading-relaxed text-muted-foreground">
          当方は、取得した個人情報等を、以下の目的のために利用します。
        </p>
        <ul className="mt-3 list-disc space-y-1 pl-6 text-muted-foreground">
          <li>本サービスの提供、運営及び本人確認のため</li>
          <li>本サービスの利用状況の分析及びサービスの維持・改善のため</li>
          <li>JukuBox API キーによる API 認証のため</li>
          <li>本サービスに関するお知らせや重要な変更等を通知するため</li>
          <li>お客様からのお問い合わせに対応するため</li>
          <li>本サービス上で著者情報及び著者が作成したコースを公開・提供するため</li>
          <li>
            お客様の利用者アカウントと著者プロフィールを紐付けて管理するため（当該紐付け情報は、お客様の利用者アカウントの削除に伴い削除します）
          </li>
          <li>
            著者として作成・公開されたコースについて、著者となったお客様が退会された後も、当該コースを購入・受講した他のお客様が継続的に閲覧・購入・受講できるよう、著者情報及びコースを保持し、公開を継続するため（この場合、お客様の退会と同時に、著者名・著者プロフィール等の表示情報は匿名化されます）
          </li>
          <li>
            Google Analytics（GA4）を利用した本サービスの利用状況の統計的な分析のため
          </li>
          <li>本ポリシーまたは本サービスの利用規約等に違反する行為への対応のため</li>
          <li>法令に基づく対応のため</li>
        </ul>

        <h2 className="mt-12 border-b border-border pb-2 font-serif text-xl md:text-2xl">
          3. 個人情報の第三者への提供
        </h2>
        <p className="mt-4 leading-relaxed text-muted-foreground">
          当方は、お客様の個人情報を、あらかじめお客様の同意を得ることなく第三者に提供することはありません。ただし、以下の場合はこの限りではありません。
        </p>
        <ul className="mt-3 list-disc space-y-1 pl-6 text-muted-foreground">
          <li>法令に基づく場合</li>
          <li>
            人の生命、身体または財産の保護のために必要がある場合であって、お客様の同意を得ることが困難であるとき
          </li>
          <li>
            国の機関もしくは地方公共団体またはその委託を受けた者が法令の定める事務を遂行することに対して協力する必要がある場合であって、お客様の同意を得ることにより当該事務の遂行に支障を及ぼすおそれがあるとき
          </li>
          <li>
            利用目的の達成に必要な範囲で、第4項に定める委託先に個人データの取扱いを委託する場合
          </li>
          <li>
            お客様が著者として登録された情報及び著者として作成・公開されたコースを、第1項(4)に定めるとおり本サービス上で公開する場合
          </li>
        </ul>

        <h2 className="mt-12 border-b border-border pb-2 font-serif text-xl md:text-2xl">
          4. 個人データの取扱いの委託
        </h2>
        <p className="mt-4 leading-relaxed text-muted-foreground">
          当方は、利用目的の達成に必要な範囲で、個人データの取扱いの全部または一部を第三者に委託することがあります。委託にあたっては、十分な個人情報保護の水準を備える者を選定し、契約等により委託先における安全管理措置が適切に講じられるよう必要かつ適切な監督を行います。本サービスにおける主な委託先は以下のとおりです。
        </p>
        <ul className="mt-3 list-disc space-y-1 pl-6 text-muted-foreground">
          <li>Google LLC（OAuth 認証基盤の提供）</li>
          <li>Google LLC（Google Analytics によるアクセス解析）</li>
          <li>
            Supabase Inc.（認証情報及びデータベースのホスティング、ファイルストレージの提供）
          </li>
          <li>Vercel Inc.（本サービスのホスティング及び配信基盤の提供）</li>
          <li>
            Amazon Web Services, Inc.（クラウドインフラストラクチャの提供。当方が直接利用するメール配送・ファイルストレージ等の各種クラウドサービス、及び Supabase Inc. が利用する基盤クラウドの提供を含みます。）
          </li>
        </ul>

        <h2 className="mt-12 border-b border-border pb-2 font-serif text-xl md:text-2xl">
          5. 外国にある第三者への提供
        </h2>
        <p className="mt-4 leading-relaxed text-muted-foreground">
          当方は、前項の委託に伴い、お客様の個人データを外国にある第三者に提供することがあります。提供先となる外国にある第三者及び当該国は以下のとおりです。
        </p>
        <ul className="mt-3 list-disc space-y-1 pl-6 text-muted-foreground">
          <li>Google LLC（アメリカ合衆国）</li>
          <li>
            Supabase Inc.（アメリカ合衆国）。なお、当方が利用する Supabase の本サービスに関するデータの保管場所は、東京リージョン（ap-northeast-1）です。
          </li>
          <li>Vercel Inc.（アメリカ合衆国）</li>
          <li>
            Amazon Web Services, Inc.（アメリカ合衆国）。なお、当方が直接利用する AWS のサービスに関するデータの保管場所は、原則として東京リージョン（ap-northeast-1）です。
          </li>
        </ul>
        <p className="mt-4 leading-relaxed text-muted-foreground">
          当該国における個人情報の保護に関する制度等の情報は、個人情報保護委員会のウェブサイト（
          https://www.ppc.go.jp/personalinfo/legal/kaiseihogohou/#gaikoku ）をご参照ください。委託先における個人情報の保護のための措置については、当該委託先との間で契約等により安全管理措置が継続的に講じられるよう必要かつ適切な監督を行います。
        </p>

        <h2 className="mt-12 border-b border-border pb-2 font-serif text-xl md:text-2xl">
          6. 安全管理措置
        </h2>
        <p className="mt-4 leading-relaxed text-muted-foreground">
          当方は、お客様の個人情報の漏えい、滅失または毀損の防止その他の安全管理のために、以下の措置を講じています。
        </p>
        <ul className="mt-3 list-disc space-y-1 pl-6 text-muted-foreground">
          <li>
            組織的安全管理措置：個人情報の取扱いに関する責任者を定め、本ポリシーの遵守状況を継続的に確認する体制を整備しています。
          </li>
          <li>
            技術的安全管理措置：通信経路を TLS により暗号化するとともに、認証情報及び API キーは適切な方式（API キーは SHA256 ハッシュ化）で保管し、平文で保存することはありません。また、API アクセスは JWT 署名検証によって認証します。
          </li>
          <li>
            物理的安全管理措置：個人データを取り扱う情報システムは、セキュリティ管理が施された信頼できるクラウド事業者の設備を利用しています。
          </li>
          <li>
            外的環境の把握：外国（アメリカ合衆国）において個人データを取り扱うにあたり、当該国の個人情報の保護に関する制度等を把握したうえで安全管理措置を講じています。
          </li>
        </ul>

        <h2 className="mt-12 border-b border-border pb-2 font-serif text-xl md:text-2xl">
          7. 個人情報等の保存期間
        </h2>
        <p className="mt-4 leading-relaxed text-muted-foreground">
          当方は、お客様の個人情報等を、利用目的の達成に必要な範囲で、以下の期間保存します。法令により保存義務がある情報については、当該法令が定める期間を上記に優先して適用します。
        </p>
        <ul className="mt-3 list-disc space-y-1 pl-6 text-muted-foreground">
          <li>
            利用者アカウント情報（ニックネーム、自己紹介文、認証情報、JukuBox API キー等）：お客様による退会のお申し出があった日、または最終ログイン日から 1 年が経過する日のいずれか早い日までに削除します。なお、不正利用への対応、係争対応その他法令上または契約上の必要がある場合には、当該対応の終了まで保存期間を延長することがあります。
          </li>
          <li>
            受講・学習履歴（コース受講履歴、トピック学習進捗等）：利用者アカウント情報と同じ期間、保存します。
          </li>
          <li>お問い合わせ履歴：対応完了から 3 年保存します。</li>
          <li>
            アクセスログ（IP アドレス、User-Agent、アクセス日時等のサーバーログ）：取得から 1 年保存します。
          </li>
          <li>
            Google Analytics により収集される情報：Google Analytics（GA4）の設定に従い、最大 14 ヶ月保存します。
          </li>
          <li>
            利用者アカウントと著者プロフィールの紐付け情報：利用者アカウント情報と同じ期間、保存します（利用者アカウントの削除に伴い削除します）。
          </li>
          <li>
            著者情報（著者名、著者プロフィール、著者スラッグ）：レコード自体は原則として無期限に保存します。ただし、お客様（著者となった利用者）の退会のタイミングで、著者名・著者プロフィール等の表示情報を当方が即時に匿名化し（例：「退会済み著者」と表示する等）、それ以降は当該著者情報から特定の個人を識別することができない状態となります。著者スラッグは、コース URL の継続性確保のために維持される場合があります。
          </li>
          <li>
            著者として作成・公開されたコース及び付随コンテンツ：原則として無期限に保存し、本サービス上での公開を継続します。
          </li>
        </ul>

        <h2 className="mt-12 border-b border-border pb-2 font-serif text-xl md:text-2xl">
          8. 保有個人データの開示等の請求手続き
        </h2>
        <p className="mt-4 leading-relaxed text-muted-foreground">
          お客様は、当方が保有するお客様自身の個人データについて、個人情報保護法に基づき、以下の請求を行うことができます。
        </p>
        <ul className="mt-3 list-disc space-y-1 pl-6 text-muted-foreground">
          <li>利用目的の通知</li>
          <li>開示（電磁的記録による方法を含む）</li>
          <li>第三者提供記録の開示</li>
          <li>内容の訂正、追加または削除</li>
          <li>利用の停止または消去</li>
          <li>第三者への提供の停止</li>
        </ul>

        <h3 className="mt-6 text-base font-semibold">請求の方法及び本人確認</h3>
        <p className="mt-2 leading-relaxed text-muted-foreground">
          請求は、本ポリシー末尾のお問い合わせ窓口までお申し出ください。当方は、本サービスにご登録のメールアドレスその他当方が指定する方法により、お客様ご本人からの請求であることを確認のうえ、合理的な期間及び範囲で対応いたします。なお、本サービスにはアカウント削除機能を用意しており、当該機能のご利用は、お客様自身による削除請求の手段として位置づけられます。
        </p>

        <h3 className="mt-6 text-base font-semibold">手数料</h3>
        <p className="mt-2 leading-relaxed text-muted-foreground">
          利用目的の通知及び開示の請求については、1 件あたり 1,000 円（税込）の手数料を申し受けます。請求受理後、当方が指定する方法によりお支払いください。手数料が支払われない場合は、当該請求に応じないことがあります。なお、訂正・追加・削除、利用停止、第三者提供停止の請求については、手数料を申し受けません。
        </p>

        <h3 className="mt-6 text-base font-semibold">
          著者情報及びコースに関する特例
        </h3>
        <p className="mt-2 leading-relaxed text-muted-foreground">
          お客様が著者として登録された情報、及び著者として作成・公開されたコースは、本サービス上で公開し継続的に提供することを目的として登録・作成されるものです。お客様の退会後も、当該コースを購入・受講した他のお客様への継続的なサービス提供のため、著者レコード及びコースは本サービス上に保持され、公開が継続されます。
        </p>
        <p className="mt-4 leading-relaxed text-muted-foreground">
          ただし、お客様の退会と同時に、当方は著者名・著者プロフィール等の表示情報を即時に匿名化（例：「退会済み著者」と表示する等）します。これにより、退会後の著者情報から特定の個人を識別することはできなくなります。退会後に改めて著者情報の削除をご請求いただく必要は原則としてありません。
        </p>
        <p className="mt-4 leading-relaxed text-muted-foreground">
          なお、お客様の退会前における著者情報の削除請求があった場合、他のお客様の購入及び受講の継続を保護するため、著者レコードそのものの削除には応じられないことがあります。この場合も、合理的な代替措置として、著者名・著者プロフィール等の表示情報の匿名化対応を行います。コース本体の内容（タイトル・説明・トピック等）については、原則として削除請求の対象外となり、その権利関係は本サービスの利用規約に従います。
        </p>

        <h2 className="mt-12 border-b border-border pb-2 font-serif text-xl md:text-2xl">
          9. Cookie 等の利用について
        </h2>
        <p className="mt-4 leading-relaxed text-muted-foreground">
          本サービスでは、以下の目的で Cookie 及び類似技術を利用しています。
        </p>

        <h3 className="mt-6 text-base font-semibold">
          (1) ログインセッション維持のための Cookie
        </h3>
        <p className="mt-2 leading-relaxed text-muted-foreground">
          お客様のログインセッションを維持するために必要な Cookie を利用しています（HttpOnly Cookie として安全に管理されます）。当方は、本 Cookie を広告配信や行動追跡の目的で利用することはありません。
        </p>

        <h3 className="mt-6 text-base font-semibold">
          (2) Google Analytics による利用状況の分析のための Cookie
        </h3>
        <p className="mt-2 leading-relaxed text-muted-foreground">
          本サービスでは、サービス改善の目的で、Google LLC が提供するアクセス解析サービス Google Analytics（GA4）を利用しています。Google Analytics は Cookie を利用してお客様の閲覧情報を取得し、Google LLC（アメリカ合衆国）に送信します。送信される情報には、Cookie 識別子（Client ID）、IP アドレス、User-Agent、閲覧 URL、リファラー、閲覧日時等が含まれ、本サービスの利用状況の統計的な分析の目的にのみ利用されます。
        </p>
        <p className="mt-4 leading-relaxed text-muted-foreground">
          Google Analytics の利用規約及びプライバシーポリシーについては、Google LLC のウェブサイトをご参照ください。Google Analytics によるデータ収集を無効にしたい場合は、Google LLC が提供するオプトアウトアドオン（
          https://tools.google.com/dlpage/gaoptout?hl=ja ）をご利用ください。
        </p>

        <h3 className="mt-6 text-base font-semibold">外部送信規律に基づく公表事項</h3>
        <p className="mt-2 leading-relaxed text-muted-foreground">
          電気通信事業法に定める外部送信規律に基づき、本サービスからお客様の端末以外の電気通信設備に送信される利用者情報は、上記 (2) に記載の Google Analytics に関する送信のみです。送信先（Google LLC、アメリカ合衆国）、送信される情報の種類、利用目的は上記のとおりです。
        </p>

        <h3 className="mt-6 text-base font-semibold">Cookie の無効化について</h3>
        <p className="mt-2 leading-relaxed text-muted-foreground">
          お客様は、ブラウザの設定により Cookie の受け入れを拒否することができます。ただし、Cookie を無効化した場合は、本サービスへのログインが行えない等、本サービスの一部または全部をご利用いただけないことがあります。
        </p>

        <h2 className="mt-12 border-b border-border pb-2 font-serif text-xl md:text-2xl">
          10. プライバシーポリシーの変更手続き
        </h2>
        <p className="mt-4 leading-relaxed text-muted-foreground">
          当方は、法令の改正、本サービスの内容の変更その他の事情により、本ポリシーを変更することがあります。変更後の内容は、本ページに掲載した時点から効力を生じるものとします。お客様にとって重要な変更を行う場合には、本サービス上で事前にお知らせします。
        </p>

        <h2 className="mt-12 border-b border-border pb-2 font-serif text-xl md:text-2xl">
          11. お問い合わせ及び苦情等の窓口について
        </h2>
        <p className="mt-4 leading-relaxed text-muted-foreground">
          本ポリシーの内容、本サービスにおける個人情報等の取扱いその他に関するお問い合わせ及び苦情等は、以下の窓口までご連絡ください。
        </p>
        <ul className="mt-3 list-disc space-y-1 pl-6 text-muted-foreground">
          <li>
            事業者：屋号「はるさめ dev」にて本サービスを運営する個人事業主です。代表者氏名については、保有個人データの開示請求等の手続きを通じてご本人に開示いたします。
          </li>
          <li>事業者の所在地：バーチャルオフィス（後日記載）</li>
          <li>お問い合わせ窓口：お問い合わせフォーム（/contacts）</li>
        </ul>

        <p className="mt-12 font-mono text-xs text-muted-foreground">
          制定日: 2026-05-04
        </p>
      </article>
    </main>
  );
}
