FROM mcr.microsoft.com/dotnet/sdk:9.0 AS build

WORKDIR /app

COPY . ./

RUN dotnet restore
#It would be safer to do it in CI/CD.
RUN dotnet test "./tests/RatingService.Tests" --no-build --verbosity normal 
RUN dotnet publish -c Release -o out

FROM mcr.microsoft.com/dotnet/aspnet:9.0 AS runtime

WORKDIR /app

COPY --from=build /app/out ./

ENTRYPOINT [ "dotnet", "RatingService.dll" ]
