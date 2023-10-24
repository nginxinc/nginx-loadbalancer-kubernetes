# Demo commands

## Creation / Updte / Deletion of ConfigMap

Run the configuration test harness:

```bash
go run cmd/configuration-test-harness/main.go
```

Run each of the following commands in turn and watch the output of the configuration test harness.

Create:

```bash
cat <<'EOF' | kubectl apply -f -
apiVersion: v1
kind: ConfigMap
metadata:
  name: nlk-config
  namespace: nlk
data:
  nginx-hosts: "http://10.1.1.4:9000/api,http://10.1.1.5:9000/api"
  tls-mode: "ss-mtls"
  ca-certificate: "nlk-tls-ca-secret"
  client-certificate: "nlk-tls-client-secret"
EOF
```

Update:

```bash
cat <<'EOF' | kubectl apply -f -
apiVersion: v1
kind: ConfigMap
metadata:
  name: nlk-config
  namespace: nlk
data:
  nginx-hosts: "http://10.1.1.4:9000/api,http://10.1.1.5:9000/api,http://10.1.1.6:9000/api"
  tls-mode: "ss-mtls"
  ca-certificate: "nlk-tls-ca-secret"
  client-certificate: "nlk-tls-client-secret"
EOF
```

Delete:

```bash
kubectl delete configmap nlk-config -n nlk
```

## Creation / Update / Deletion of Certificate Secrets

Run the certificates test harness:

```bash
go run cmd/certificates-test-harness/main.go
```

Run each of the following commands in turn and watch the output of the certificates test harness.

Create:

```bash
cat <<'EOF' | kubectl apply -f -
apiVersion: v1
data:
  tls.crt: LS0tLS1CRUdJTiBDRVJUSUZJQ0FURS0tLS0tCk1JSURUekNDQWpjQ0ZGNGlKWWxnWFYrNksxclE0L1JacXEyTWxRTWZNQTBHQ1NxR1NJYjNEUUVCQ3dVQU1HUXgKQ3pBSkJnTlZCQVlUQWxWVE1STXdFUVlEVlFRSURBcFhZWE5vYVc1bmRHOXVNUkF3RGdZRFZRUUhEQWRUWldGMApkR3hsTVE0d0RBWURWUVFLREFWT1IwbE9XREVlTUJ3R0ExVUVDd3dWUTI5dGJYVnVhWFI1SUNZZ1FXeHNhV0Z1ClkyVnpNQjRYRFRJek1UQXdNakl6TURReU9Wb1hEVEl6TVRFd01USXpNRFF5T1Zvd1pERUxNQWtHQTFVRUJoTUMKVlZNeEV6QVJCZ05WQkFnTUNsZGhjMmhwYm1kMGIyNHhFREFPQmdOVkJBY01CMU5sWVhSMGJHVXhEakFNQmdOVgpCQW9NQlU1SFNVNVlNUjR3SEFZRFZRUUxEQlZEYjIxdGRXNXBkSGtnSmlCQmJHeHBZVzVqWlhNd2dnRWlNQTBHCkNTcUdTSWIzRFFFQkFRVUFBNElCRHdBd2dnRUtBb0lCQVFEc0xteW9vc2NKNXZQbmFER3FCZkozYmQxT2F0NXYKYjQwSjcrc0ROamxhUFNGRjdNZGJyN2JFOFhkRUNGOHZCYkE2MkUwTVJRYXJHMFl5aWJBOG05OFV5TVJ1R0VRTApSZkltZ2pHVWtXdkRmU0ViSkJLZ1RXay83ckJPeDVRY1lFVnlCZ2Z2SDZMWk5hampic0VFaUFjalVObkp3WkNpCjdIWjJDSVZ5NVhpUmlKWVZLbGZwN3QvYlIwTXRtWVNEMlh6L0RoZXBRMnpkQ25zNnpFK3Q5aGxINnorcThydUUKam45ZTJaUGNBeTlHc1lNVi9WSTdIK0ZyMVBaREJJNUhtb1pZeUo4MTVqcWhiQUxGeU00Z0x3V1JlSVJuek9aZApNeVV0QW8rM1Vic1ZwQ0NOTE9XSk1WeXVUUVV1U3QwOFlXNTFDeGtESlBUMDY0bFAyZXBLQ2RQcEFnTUJBQUV3CkRRWUpLb1pJaHZjTkFRRUxCUUFEZ2dFQkFMUStGRzVxNDVSV0JFK1JwNkpwUlZWakphMFF6ZE44V1ljcWUwTDkKVlJPbG5aSkxva2JGNzJrOEdrZ0t3YzNVZ21IRzFMbmMvdzlTbFJ1Ti9uRkMySEtZbWxqWHkwaHpWaFZLbDRjRgozc1NKV0c2dzBkcXFNVmk0bDBybFVDQlk3cDBTQXA4eFdLd2ZoanE0NU1YSGlwNERRUzVTUXUrWjJiSCtTN3pzCjNxcGhDaEpkNkpGZUNnL0pGQ1VIWjZGRXk0ZitJcGl1bU01SENKakJUcHdEVUdLcHB3a2pPeWp6Ym1Vc1hEcUYKRkIvYzZ4MnluYjlROGQ4MG1oREV5V0dFWDhtbE1ycXNrblFZd2JTenkyUitudnFLOFp3NTFGbDRSUU96K1lPZAphMUlTemZKWkYwdEtqcittS0ZSODZXWm05VWNOT1JnWjQ1WjVRNWI1V1AvMVF6dz0KLS0tLS1FTkQgQ0VSVElGSUNBVEUtLS0tLQo=
  tls.key: LS0tLS1CRUdJTiBQUklWQVRFIEtFWS0tLS0tCk1JSUV2QUlCQURBTkJna3Foa2lHOXcwQkFRRUZBQVNDQktZd2dnU2lBZ0VBQW9JQkFRRHNMbXlvb3NjSjV2UG4KYURHcUJmSjNiZDFPYXQ1dmI0MEo3K3NETmpsYVBTRkY3TWRicjdiRThYZEVDRjh2QmJBNjJFME1SUWFyRzBZeQppYkE4bTk4VXlNUnVHRVFMUmZJbWdqR1VrV3ZEZlNFYkpCS2dUV2svN3JCT3g1UWNZRVZ5Qmdmdkg2TFpOYWpqCmJzRUVpQWNqVU5uSndaQ2k3SFoyQ0lWeTVYaVJpSllWS2xmcDd0L2JSME10bVlTRDJYei9EaGVwUTJ6ZENuczYKekUrdDlobEg2eitxOHJ1RWpuOWUyWlBjQXk5R3NZTVYvVkk3SCtGcjFQWkRCSTVIbW9aWXlKODE1anFoYkFMRgp5TTRnTHdXUmVJUm56T1pkTXlVdEFvKzNVYnNWcENDTkxPV0pNVnl1VFFVdVN0MDhZVzUxQ3hrREpQVDA2NGxQCjJlcEtDZFBwQWdNQkFBRUNnZ0VBWWV3Tm1RMkVRSkxFNVZqSjNwaUFvd3VtQ2ZFOU1DNnI1MGJWeFlzaDFFd3MKRTNYTVlqTkVML3Q5VzNPdEl5M1VsMUUvQUt0TnpIdU9hejJ6R0MzNEhBSHhqMFA0VWtRNTFjVjlFUUFLRWc4NwpQcW1DSDN4NCtzelh4Skh5MHFFSHFmTGVMMEtLbmt3bExjYXB1RnM5dW1LM0tYTmJxSEVwM0Y1RUZoTVdIaUFiClV4cWpYZXJKWXBZMXU3RVhwUHFWaSt3dFhkWDNuT092bStZMzNLemRGNTErTXBpOHVkMGpTY0p2VEtIYXpKV1oKZE40aFF4SEpEOGpMZm8zNDRNWEFYQndkWFVuN1V4NW9qU0dVemtnU1pKTnBkQVNaUU5lN2E1Qk4xalREZ29LaQp0My9kS0NFTmFsMGdUQVI1MFArYWJxeWhSOFUxbVVuT1dsWDN1ekJsa1FLQmdRRDNFS09tWm16Z293Nk5mOWp2CmpkdzJOUUQ1Q2NlWEk4cWlxRjlMNWZoVWV0b29PajA2WlZpYjBEejNWcWlqVnRITFBQV2pRQ2ZwWC9DVHdkUTkKajNHTVBDOEhIZWJkeEZqTHpRSWlDUXAvaGFlc2RGTUhHRzFBak92VDRNdkxOaklWM29JOCszZjhIWk81Q0pYQgpJdW1xcG82YWJhV094VEhDS2pYY2YxNTI0d0tCZ1FEMHVRWGNJWlJrNjdrYU5PZWhjLzFuU2QyTVJSVG1zV3JYClFjTEtQczdWVFkwWlA1T2hsU1pRVkM2MW9PUWtuZnFScGhVbnJVVDF1NXFRbzJUSGhGcmN4a3Qrd1hCenhDWDgKUXZXY1RTTmkvZ0Y2Rkx1OEsyYXlOVlBQTDFNQjJ4ZGxTL0Z0SDRLSG5RWHM1NFM5MVZHRTR5OUd6aDIwWEcwOQp4K0FiQzlPM3d3S0JnSGNuWURXbFdrY3dmSmxEbW1WV0htbEtRTkRhcFphLzNUOTdRcEtCTTdYU2xob21sRmJ3CmY3Nk52SWx4RXQzTHhseGxadlkzdjhmdXpFRUdqd3l0Zkk2c2krVzd4eGNYVmRmY1pIWHp0RXR5TXo2WnoxMHgKcTZjaEQ2OWMwQXlPYzdOV1g2dDNnQk5vVkZFOTBiT1cyZWpDY1M0TFNYaEVwRTNITzdpKytOa1BBb0dBZW9CYgo1Sk9TbXVvOG9GZTNVMlNpaHArOUhVZy9iRE9IamZWSE1zSTUreUIwN3h5YUpCcHJNVzdTYXV6OUJ5OWxqSjhjCm05M3FWUy94OFZFNVUzNTNsV2hWeGovQ3NOQ1JTek9oaXZvNktvV0g2N3FSTjJKcVorNjE0MUtITkxpZGY0R0MKZXVONURiV1dqNzVjL2tIWUtyTW1xVVRvTGE3T3FFeHpiRmFCUnMwQ2dZQUN6UlBiQTBqK0JyOWpxenFKYU56VQpPZzJ4aUFhZlNuTzBxTkREZWZNdXhZSnFMZ0llMG9HRUExb1JlcDhvTDNDeWZYd2h3U2xqMlRjeGZZTFR5MENKCmZBejF4QjY3SUF3cGlvZnI0K1ZXekF4SC9BbnlhVGNZVGdtakdKdlhtZ3l1RGU3ZGJsa3lJU1M5ZTBLc1JzeTYKQ2hPdVZoQy96bjErbm9BMkNsdzFJdz09Ci0tLS0tRU5EIFBSSVZBVEUgS0VZLS0tLS0K
kind: Secret
metadata:
  creationTimestamp: null
  name: nlk-tls-ca-secret
  namespace: nlk
type: kubernetes.io/tls
EOF
```

Update:

```bash
cat <<'EOF' | kubectl apply -f -
apiVersion: v1
data:
  tls.crt: LS0tLS1CRUdJTiBDRVJUSUZJQ0FURS0tLS0tCk1JSUVsekNDQTMrZ0F3SUJBZ0lVUFRrVExlS3FNQlRlNFZhQXpiL1hXcSt6THRBd0RRWUpLb1pJaHZjTkFRRUwKQlFBd1pERUxNQWtHQTFVRUJoTUNWVk14RXpBUkJnTlZCQWdNQ2xkaGMyaHBibWQwYjI0eEVEQU9CZ05WQkFjTQpCMU5sWVhSMGJHVXhEakFNQmdOVkJBb01CVTVIU1U1WU1SNHdIQVlEVlFRTERCVkRiMjF0ZFc1cGRIa2dKaUJCCmJHeHBZVzVqWlhNd0hoY05Nak14TURBeU1qTXdPRFEyV2hjTk1qUXhNREF4TWpNd09EUTJXakI3TVFzd0NRWUQKVlFRR0V3SlZVekVUTUJFR0ExVUVDQXdLVjJGemFHbHVaM1J2YmpFUU1BNEdBMVVFQnd3SFUyVmhkSFJzWlRFTwpNQXdHQTFVRUNnd0ZUa2RKVGxneEhqQWNCZ05WQkFzTUZVTnZiVzExYm1sMGVTQW1JRUZzYkdsaGJtTmxjekVWCk1CTUdBMVVFQXd3TWJYbGtiMjFoYVc0dVkyOXRNSUlCSWpBTkJna3Foa2lHOXcwQkFRRUZBQU9DQVE4QU1JSUIKQ2dLQ0FRRUF2RjJSVklsVEQvQWVSQkZSNTNrb1dlZmlhQ09keVJvWnJRSm04RGZJdk5vQ2pTdkNPUEhzZ2VSSQptcmhxcEhkc2FSK2RqSjNRbzFKemljeS9YNHgxRy9SSGRNM0d3OXNuTS9jdHBvb1FSZUtFTi90T1FnVnBMaEZQClE4K3cvaDVRSGZWWFFQcFdyeGhjTkFiK1hrbzk5SVBndE81blNvazExK2pEOGhJSVlTTUEvVElwUFBPeXNKdnkKOUhKUjBEY3pqWE1VQTR0dGphYTZ0cm8xNm81WEZsdDZPM251YTFkc0VLYVRJQmhWbXI4ZGJWVFd2c1VNYTV6WgpyR0pSa25Bc01RV1lxSVBBVG1MdXlaeHlYTG1WQkp4U1FzWTZqYVllLzcreXZGdUVCUTZ4MDByZ1Z0THpzWi9NCitZSHg2VUVXRFoxYzZuZWxxdlFja1FvSzdTd3VLUUlEQVFBQm80SUJLRENDQVNRd0NRWURWUjBUQkFJd0FEQUwKQmdOVkhROEVCQU1DQmVBd1NBWURWUjBSQkVFd1A0SU1iWGxrYjIxaGFXNHVZMjl0Z2hOelpYSjJaWEl1YlhsawpiMjFoYVc0dVkyOXRnZzRxTG0xNVpHOXRZV2x1TG1OdmJZY0VDZ0FBQ29jRUNnQUFDekFUQmdOVkhTVUVEREFLCkJnZ3JCZ0VGQlFjREFUQWRCZ05WSFE0RUZnUVVkOUVieTR0TDNZSjhqSTd6RzVKR2p5TW5vVDR3Z1lzR0ExVWQKSXdTQmd6Q0JnS0ZvcEdZd1pERUxNQWtHQTFVRUJoTUNWVk14RXpBUkJnTlZCQWdNQ2xkaGMyaHBibWQwYjI0eApFREFPQmdOVkJBY01CMU5sWVhSMGJHVXhEakFNQmdOVkJBb01CVTVIU1U1WU1SNHdIQVlEVlFRTERCVkRiMjF0CmRXNXBkSGtnSmlCQmJHeHBZVzVqWlhPQ0ZGNGlKWWxnWFYrNksxclE0L1JacXEyTWxRTWZNQTBHQ1NxR1NJYjMKRFFFQkN3VUFBNElCQVFCRFY0OWtOanpCZmpPMjlNWVd1ZFlMNktzYlFrdlFDdzZ3dEF6cVJ2Smd5WEtBbXYzcwp6OFJjZHFNRTZ2bXdWZUxpS1ZDeDFSOXY1dlhBS0hMNmFSSi9QMk1QUm1Ic0pNM2MxSVF2NjAzVVYzaCtQL3JICll3QTVVRzBCMXlMaWs3N3RIdWZtRGFKdlRZMGpvcWJ0cVEwU1ByWjh6TGJQdmo3VTc1bStHWEx5VEp0UnRjdDUKSVJTblg0NlZxbWs0Q1l2dW1SaDJJTy8yUHpKNkpUbzA5VFpEOGpDbyttbHIzSXhXOHBsRDBYSkQ1YXY4cEQ2dApDZ2xnR2FyakRWSmVCQUZobkVIWStYelZDYzJVSktVblNiaDVWUVRLdno2dCtCWTlldi9LOFJGYlcxWVd3N1Q1CmQ0UzEyaC94cUpXcUJOdXFJeldkMTMvbTZaaHoraHorb0VvRQotLS0tLUVORCBDRVJUSUZJQ0FURS0tLS0tCg==
  tls.key: LS0tLS1CRUdJTiBQUklWQVRFIEtFWS0tLS0tCk1JSUV2UUlCQURBTkJna3Foa2lHOXcwQkFRRUZBQVNDQktjd2dnU2pBZ0VBQW9JQkFRQzhYWkZVaVZNUDhCNUUKRVZIbmVTaFo1K0pvSTUzSkdobXRBbWJ3TjhpODJnS05LOEk0OGV5QjVFaWF1R3FrZDJ4cEg1Mk1uZENqVW5PSgp6TDlmakhVYjlFZDB6Y2JEMnljejl5Mm1paEJGNG9RMyswNUNCV2t1RVU5RHo3RCtIbEFkOVZkQStsYXZHRncwCkJ2NWVTajMwZytDMDdtZEtpVFhYNk1QeUVnaGhJd0Q5TWlrODg3S3dtL0wwY2xIUU56T05jeFFEaTIyTnBycTIKdWpYcWpsY1dXM283ZWU1clYyd1FwcE1nR0ZXYXZ4MXRWTmEreFF4cm5ObXNZbEdTY0N3eEJaaW9nOEJPWXU3SgpuSEpjdVpVRW5GSkN4anFOcGg3L3Y3SzhXNFFGRHJIVFN1Qlcwdk94bjh6NWdmSHBRUllOblZ6cWQ2V3E5QnlSCkNncnRMQzRwQWdNQkFBRUNnZ0VBQnRtRGh6dnhKeU5pV1FqNk0xSjVaYW81Z1BPZW5CUC9aVnV3MWtFVHVIa0QKQ1JMVGRCSlVEb3NqR3NFNXNONWVOUVVVZG1zU0RiRUVzNDVHL1VOVWRONUdyNWdyWEs3bzk0cExBTmlFd1UzUgo1TFBybU1TdFJQZ1YxaytaK09aWVNudnVudFhHUzRWZ0ZrenFJMlZNRXlBQ2pxeHhXWFFHTkdJL1VkdXNXS3FICkNlOWNGMXp0UGxBWkVIRU9JNTJEcVFZK3FhUUtNSjVIT3RNS1F6SzY5MGF0VUJmbWZxUCsxYTh4Vm1VV1cyaVgKNWpuN3IzeDRJdENwSlpwTHM2SzZRS29aNU9qd25nSTc4MkxuQ3V0cC9NZmg5RmVYSlFRRVppYXMzOW44c29NcwpCN0ZKQnpKalQwVlQxR2hVWnNwSUhra1AwR3pBQTd2Y3MwNitDQWc5Z1FLQmdRRG1CYTQ3MnlLYmwvZUIyY3BFCjdZeXVmaHk1M093OTFzNkkwYnQxYWZyajNvcHExV1VVWUJiOWhIUm1CRXkzVVE2V1JoenltdG1MdXBRTWh4enkKeThSMVFtNW1RUTI1T3lsNGl3TFZQdVE3Ty94dzh6UUdNR1pBY08zQ3NuYUM5Y3YxZFpnZzcvQ1BUTlpsY2RxZAp2OXFqNnpwVHFrT21OVnlFUGxPeFI3aHJ5UUtCZ1FEUm80WVNwcDRONHRlSWhIaG1uZGpYaVZkSDYyM3oxQldICnV1U3M0L0M2RDVjaUhWRTJYZnZ2cnQ4eEhBbEVycGtoaVFUZ1JzZHVwblVhVVFodUtaT1RNTXVGcDAzblA1cmMKUjlxY1MxQUVQUXZMVFNZN0dqaEptaExEVUM0WktGYmo4bXhkSWVBK09YL21mank4SFJTMWYxYncwR3owZCtvNQozVnlyOUU0ZllRS0JnUUNuY1QwckwxTGJCdDNTZFpMcmFDMC9uR2dXMkg1VWFia0JHZ09tN2hZSHFLa0VLZ0VoCnV1MGhjVGsyUml6K1NSQWdUanVtVXhqSHdYTWlSM3pJTlpMMmRQeGVqVDZMTjBqeUNlZHZDaEFrR24raVRUZnkKeFdxNXdEc2p2cnZNaTFjRWdLelVWVFc5YXdhcTVCMXJOZ3pYeEZVNk1EaDhsbDJabXJGYjNNU2dHUUtCZ0F1cQpmT2lHeXg3Y3M3L09GMkVtZ1kyay8rMXBwWW0vRUorbi85ZTdLNGMvSE5yeUpMWFF6eGRNZFBFbnJVQmNNdnRSCnc2cXpaWis3dGFLTVJkclRoM25XYWt6NnZYUVQ3d3M1R0dwQUtxakJ1T2xNVnNkTk16cXRUMFA5TDBPSklpUzMKTmQ2TTV3eXZhSFdzS3JjUkt6amFhRDBvYkJmQ29JOHR5VjFzVC9paEFvR0FWNHJzMU14d2VkRVVPTFo4Y0ZBQgpvYTE1MG41dVgwUGQ1Nm92enRXNHZoR2JyOVpjQ294dnR0a0VqQUNvTzdyanhBb1JyYVk1U1A0WTdHaHJxb043CkZYeFYrVVh0aGxyU21SQWUzS25ndzA4bkRpT0Fia0NZdFlGVGs3d3d3V0xYU2VYUWRCdUx1aWRqZ05PL2RhbkkKczUzVXU3OGlnRVdPa21YWThGYm1OVEU9Ci0tLS0tRU5EIFBSSVZBVEUgS0VZLS0tLS0K
kind: Secret
metadata:
  creationTimestamp: null
  name: nlk-tls-ca-secret
  namespace: nlk
type: kubernetes.io/tls
EOF
```

Delete:

```bash
kubectl delete secret nlk-tls-ca-secret -n nlk
```

## Generation of TLS Configuration based on ConfigMap and Certificate Secrets

Run the TLS configuration test harness:

```bash
go run cmd/tls-configuration-test-harness/main.go
```

The test harness will iterate through the tls modes and show basic information about the generated TLS configuration.